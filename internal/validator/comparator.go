package validator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
	"github.com/googleapis/gapic-config-validator/internal/config"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/longrunning"
)

var (
	wellKnownPatterns = map[string]bool{
		"projects/{project}":                      true,
		"organizations/{organization}":            true,
		"folders/{folder}":                        true,
		"projects/{project}/locations/{location}": true,
		"billingAccounts/{billing_account_id}":    true,
	}

	wellKnownNames = map[string]bool{
		"project":      true,
		"organization": true,
		"folder":       true,
		"location":     true,
	}

	isIamMethod = map[string]bool{
		"GetIamPolicy":       true,
		"SetIamPolicy":       true,
		"TestIamPermissions": true,
	}
)

func (v *validator) compare() {
	// compare interfaces
	v.compareServices()

	// compare resource references
	v.compareResourceRefs()
}

func (v *validator) compareServices() {
	for _, inter := range v.gapic.GetInterfaces() {
		serv := v.resolveServiceByName(inter.GetName())
		if serv == nil {
			v.addError("Interface %q does not exist", inter.GetName())
			continue
		}

		// compare resources
		v.compareResources(inter)

		// compare methods
		for _, method := range inter.GetMethods() {
			methodDesc := serv.FindMethodByName(method.GetName())
			if methodDesc == nil {
				v.addError("Method %q does not exist", inter.GetName()+"."+method.GetName())
				continue
			}

			v.compareMethod(methodDesc, method)
		}
	}
}

func (v *validator) compareMethod(methodDesc *desc.MethodDescriptor, method *config.MethodConfigProto) {
	fqn := methodDesc.GetFullyQualifiedName()
	mOpts := methodDesc.GetMethodOptions()

	// ignore IAM methods
	if isIamMethod[method.GetName()] {
		return
	}

	// compare method_signatures & flattening groups
	if flattenings := method.GetFlattening(); flattenings != nil {
		eSigs, err := ext(mOpts, annotations.E_MethodSignature)
		if err != nil {
			v.addError("Method %q missing method_signature(s) for flattening(s)", fqn)
			goto LRO
		}
		sigs := eSigs.([]string)

		for _, flat := range flattenings.GetGroups() {
			joined := strings.Join(flat.GetParameters(), ",")
			if !containStr(sigs, joined) {
				v.addError("Method %q missing method_signature for flattening %q", fqn, joined)
			}
		}
	}

LRO:
	// compare operation_info & longrunning config
	if lro := method.GetLongRunning(); lro != nil {
		eLRO, err := ext(mOpts, longrunning.E_OperationInfo)
		if err != nil {
			v.addError("Method %q missing longrunning.operation_info", fqn)
			goto Behavior
		}
		info := eLRO.(*longrunning.OperationInfo)

		// trim to local message name
		protoRes := info.GetResponseType()
		if strings.Contains(protoRes, ".") {
			protoRes = protoRes[strings.LastIndex(protoRes, ".")+1:]
		}
		gapicRes := lro.GetReturnType()
		if strings.Contains(gapicRes, ".") {
			gapicRes = gapicRes[strings.LastIndex(gapicRes, ".")+1:]
		}

		if protoRes != gapicRes {
			v.addError("Method %q operation_info.response_type %q does not match %q",
				fqn,
				protoRes,
				gapicRes)
		}

		// trim to local message name
		protoMeta := info.GetMetadataType()
		if strings.Contains(protoMeta, ".") {
			protoMeta = protoMeta[strings.LastIndex(protoMeta, ".")+1:]
		}
		gapicMeta := lro.GetMetadataType()
		if strings.Contains(gapicMeta, ".") {
			gapicMeta = gapicMeta[strings.LastIndex(gapicMeta, ".")+1:]
		}

		if protoMeta != gapicMeta {
			v.addError("Method %q operation_info.metadata_type %q does not match %q",
				fqn,
				protoMeta,
				gapicMeta)
		}
	}

Behavior:
	// compare input message field_behaviors & required_fields
	if reqs := method.GetRequiredFields(); len(reqs) > 0 {
		input := methodDesc.GetInputType()

		for _, name := range reqs {
			field := input.FindFieldByName(name)
			if field == nil {
				v.addError("Field %q in method %q required_fields does not exist in %q",
					name,
					fqn,
					input.GetFullyQualifiedName())
				continue
			}

			eBehv, err := ext(field.GetFieldOptions(), annotations.E_FieldBehavior)
			if err != nil {
				v.addError("Field %q is missing field_behavior = REQUIRED per required_fields config", field.GetFullyQualifiedName())
				continue
			}
			behavior := eBehv.([]annotations.FieldBehavior)

			if !containBehavior(behavior, annotations.FieldBehavior_REQUIRED) {
				v.addError("Field %q is not annotated as REQUIRED per required_fields config", field.GetFullyQualifiedName())
			}
		}
	}
}

func (v *validator) compareResources(inter *config.InterfaceConfigProto) {
	for _, res := range inter.GetCollections() {
		ent := snakeToCamel(res.GetEntityName())
		pat := res.GetNamePattern()
		if wellKnownPatterns[pat] || wellKnownNames[ent] {
			continue
		}

		for _, f := range v.files {
			// check file for resource_definition annotations
			eResDef, err := ext(f.GetFileOptions(), annotations.E_ResourceDefinition)
			if err == nil {
				resDefs := eResDef.([]*annotations.ResourceDescriptor)
				for _, res := range resDefs {
					typ := res.GetType()
					typ = typ[strings.Index(typ, "/")+1:]

					if ent == typ {
						if !containStr(res.GetPattern(), pat) {
							v.addError("resource definition for %q in %q does not have pattern %q",
								res.GetType(),
								f.GetFullyQualifiedName(),
								pat)
						}
						goto NextCollectionItem
					}
				}
			}

			for _, m := range f.GetMessageTypes() {
				eRes, err := ext(m.GetMessageOptions(), annotations.E_Resource)
				if err != nil {
					continue
				}
				resDesc := eRes.(*annotations.ResourceDescriptor)

				typ := resDesc.GetType()
				typ = typ[strings.Index(typ, "/")+1:]

				// the pattern is defined in a resource named differently than the
				// name_pattern value, which is OK.
				if containStr(resDesc.GetPattern(), pat) {
					goto NextCollectionItem
				}

				if typ == ent {
					if !containStr(resDesc.GetPattern(), pat) {
						v.addError("resource definition for %q in %q does not have pattern %q",
							resDesc.GetType(),
							m.GetFullyQualifiedName(),
							pat)
					}

					goto NextCollectionItem
				}
			}
		}

		v.addError("No corresponding resource definition for %q: %q", res.GetEntityName(), res.GetNamePattern())

	NextCollectionItem:
	}
}

func (v *validator) compareResourceRefs() {
	for _, ref := range v.gapic.GetResourceNameGeneration() {
		// skip IAM messages
		if strings.HasPrefix(ref.GetMessageName(), "google.iam.v1") {
			continue
		}

		msgDesc := v.resolveMsgByLocalName(ref.GetMessageName())
		if msgDesc == nil {
			v.addError("Message %q in resource_name_generation item does not exist", ref.GetMessageName())
			continue
		}

		for fname, ref := range ref.GetFieldEntityMap() {
			// skip nested fields, presumably they are
			// being validated in the origial msg
			if strings.Contains(fname, ".") {
				continue
			}

			field := msgDesc.FindFieldByName(fname)
			if field == nil {
				v.addError("Field %q does not exist on message %q per resource_name_generation item", fname, msgDesc.GetFullyQualifiedName())
				continue
			}

			var typ, child string
			if eResRef, err := ext(field.GetFieldOptions(), annotations.E_ResourceReference); err == nil {
				resRef := eResRef.(*annotations.ResourceReference)

				typ = resRef.GetType()
				child = resRef.GetChildType()
			} else if eRes, err := ext(msgDesc.GetMessageOptions(), annotations.E_Resource); err == nil {
				res := eRes.(*annotations.ResourceDescriptor)
				typ = res.GetType()
			} else {
				v.addError("Field %q missing resource_reference to %q", field.GetFullyQualifiedName(), ref)
				continue
			}

			// use child_type instead
			if typ == "" {
				childMsg := v.resolveResRefMessage(child, msgDesc.GetFile())
				if childMsg == nil {
					v.addError("child_type %q on %q is not a defined resource", child, field.GetFullyQualifiedName())
					continue
				}

				refItem := v.resolveRefFromCollections(ref)
				if refItem == nil {
					v.addError("entity_name %q is not a defined in any available collection", ref)
				}

				if eResDef, err := ext(childMsg.GetMessageOptions(), annotations.E_Resource); err == nil {
					var found bool
					resDef := eResDef.(*annotations.ResourceDescriptor)
					for _, pattern := range resDef.GetPattern() {
						if strings.HasPrefix(pattern, refItem.GetNamePattern()) {
							found = true
							break
						}
					}

					if !found {
						v.addError("Field %q child_type %q isn't a proper child of %q in GAPIC config", field.GetFullyQualifiedName(), child, ref)
					}
				}

				continue
			}

			// compare using upper camel case names
			t := typ[strings.Index(typ, "/")+1:]
			if !wellKnownTypes[typ] && t != snakeToCamel(ref) {
				v.addError("Field %q resource_type_kind %q doesn't match %q in config", field.GetFullyQualifiedName(), typ, ref)
			}
		}
	}
}

func containBehavior(arr []annotations.FieldBehavior, behv annotations.FieldBehavior) bool {
	for _, b := range arr {
		if b == behv {
			return true
		}
	}

	return false
}

func containStr(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}

	return false
}

func (v *validator) resolveServiceByName(name string) *desc.ServiceDescriptor {
	for _, f := range v.files {
		if s := f.FindService(name); s != nil {
			return s
		}
	}

	return nil
}

func (v *validator) resolveMsgByLocalName(name string) *desc.MessageDescriptor {
	for _, f := range v.files {
		fqn := f.GetPackage() + "." + name

		if m := f.FindMessage(fqn); m != nil {
			return m
		}
	}

	return nil
}

func (v *validator) resolveRefFromCollections(name string) *config.CollectionConfigProto {
	for _, item := range v.gapic.GetCollections() {
		if item.GetEntityName() == name {
			return item
		}
	}

	for _, inter := range v.gapic.GetInterfaces() {
		for _, item := range inter.GetCollections() {
			if item.GetEntityName() == name {
				return item
			}
		}
	}

	return nil
}

func (v *validator) parseParameters(p string) error {
	for _, s := range strings.Split(p, ",") {
		if e := strings.IndexByte(s, '='); e > 0 {
			switch s[:e] {
			case "gapic-yaml":

				f, err := ioutil.ReadFile(s[e+1:])
				if err != nil {
					return fmt.Errorf("error reading gapic config: %v", err)
				}

				// throw away the first line containing
				// "type: com.google.api.codegen.ConfigProto" because
				// that's not in the proto, causing an unmarshal error
				data := bytes.NewBuffer(f)
				data.ReadString('\n')

				j, err := yaml.YAMLToJSON(data.Bytes())
				if err != nil {
					return fmt.Errorf("error decoding gapic config: %v", err)
				}

				v.gapic = &config.ConfigProto{}
				err = jsonpb.Unmarshal(bytes.NewBuffer(j), v.gapic)
				if err != nil {
					return fmt.Errorf("error decoding gapic config: %v", err)
				}
			}
		}
	}

	return nil
}

// converts snake_case and SNAKE_CASE to CamelCase.
//
// copied from github.com/googleapis/gapic-generator-go
func snakeToCamel(s string) string {
	var sb strings.Builder
	up := true
	for _, r := range s {
		if r == '_' {
			up = true
		} else if up {
			sb.WriteRune(unicode.ToUpper(r))
			up = false
		} else {
			sb.WriteRune(unicode.ToLower(r))
		}
	}
	return sb.String()
}
