package Report

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/tufin/oasdiff/diff"
)

var IsBreaking bool = false

type Report struct {
	Writer io.Writer
	level  int
}

func NewReport(buf *bytes.Buffer) *Report {
	return &Report{
		Writer: buf,
	}
}

func (r *Report) indent() *Report {
	return &Report{
		Writer: r.Writer,
		level:  r.level + 1,
	}
}

func (r *Report) print(output ...interface{}) (n int, err error) {
	return fmt.Fprintln(r.Writer, addPrefix(r.level, output)...)
}

func addPrefix(level int, output []interface{}) []interface{} {
	return append(getPrefix(level), output...)
}

func getPrefix(level int) []interface{} {
	if level == 1 {
		return []interface{}{"-"}
	}

	if level > 1 {
		return []interface{}{strings.Repeat("  ", level-1) + "-"}
	}

	return []interface{}{}
}

func (r *Report) GetTextReportAsString(d *diff.Diff, choice string) {
	switch choice {
	case "OpenApi":
		IsBreaking = false
		r.printValue(d.OpenAPIDiff, "OpenAPI")
		break
	case "Info":
		IsBreaking = false
		r.print("Info Changed")
		r.printInfoDiff(d.InfoDiff)
		break
	case "EndpointsAdded":
		IsBreaking = false
		r.print("Endpoints Added")
		r.printEndpointsAdded(d.EndpointsDiff)
		break
	// case "EndpointsDeleted":
	// 	IsBreaking = false
	// 	r.print("Endpoints Deleted")
	// 	r.printEndpointsDeleted(d)
	// 	break
	// case "EndpointsModified":
	// 	IsBreaking = false
	// 	r.print("Endpoints Modified")
	// 	r.PrintEndpointsModified(d.EndpointsDiff)
	// 	break
	case "SecurityAdded":
		IsBreaking = false
		r.print("Security Added")
		r.printSecurityAdded(d.SecurityDiff)
		break
	case "SecurityDeleted":
		IsBreaking = false
		r.print("Security Deleted")
		r.printSecurityDeleted(d.SecurityDiff)
		break
	case "SecurityModified":
		IsBreaking = false
		r.print("Security Modified")
		r.printSecurityModified(d.SecurityDiff)
		break
	case "ServerAdded":
		IsBreaking = false
		r.print("Server Added")
		r.printServerAdded(d.ServersDiff)
		break
	case "ServerDeleted":
		IsBreaking = false
		r.print("Server Deleted")
		r.printServerDeleted(d.ServersDiff)
		break
	case "ServerModified":
		IsBreaking = false
		r.print("Server Modified")
		r.printServerModified(d.ServersDiff)
		break
	case "TagAdded":
		IsBreaking = false
		r.print("Tag Added")
		r.printTagAdded(d.TagsDiff)
		break
	case "TagDeleted":
		IsBreaking = false
		r.print("Tags Deleted")
		r.printTagDeleted(d.TagsDiff)
		break
	case "TagModified":
		IsBreaking = false
		r.print("Tag Modified")
		r.printTagModified(d.TagsDiff)
		break
	case "ExternalDocs":
		IsBreaking = false
		r.print("External Docs Changed")
		r.indent().printExternalDocsDiff(d.ExternalDocsDiff)
		r.print("")
		break
	default:
		IsBreaking = false
		r.print("No changes")
	}
}
func (r *Report) printExternalDocsDiff(d *diff.ExternalDocsDiff) {
	// Added           bool
	// Deleted         bool
	// ExtensionsDiff  *ExtensionsDiff
	// DescriptionDiff *ValueDiff
	// URLDiff         *ValueDiff
	if d.Empty() {
		return
	}
	if d.Added {
		r.printConditional(d.Added, " External Docs Added")
	}
	if d.Deleted {
		r.printConditional(d.Deleted, " External Docs Deleted")
	}
	if !d.ExtensionsDiff.Empty() {
		r.printExtensionsDiff(d.ExtensionsDiff)
	}
	if !d.DescriptionDiff.Empty() {
		r.printValue(d.DescriptionDiff, "External Docs description")
	}
	if !d.URLDiff.Empty() {
		r.printValue(d.URLDiff, "External Docs URL")
	}
}
func (r *Report) printExtensionsDiff(d *diff.ExtensionsDiff) {
	if d.Empty() {
		return
	}
	sort.Sort(d.Added)
	for _, added := range d.Added {
		r.print("New Extensions Added:", added)
	}

	sort.Sort(d.Deleted)
	for _, deleted := range d.Deleted {
		r.print("Extensions Deleted", deleted)
	}
	//problem with iterarting
	keys := DiffToStringList1(d.Modified)
	sort.Sort(keys)
	for _, extension := range keys {
		r.print("Modified extensions:", extension)
		r.indent().printValue(d.Modified[extension], "Extension")
	}
}
func (r *Report) printPathDiff(d *diff.PathsDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Added)
	for _, added := range d.Added {
		r.print("New paths Added:", added)
	}

	sort.Sort(d.Deleted)
	for _, deleted := range d.Deleted {
		r.print("Paths Deleted", deleted)
	}
	//problem with iterating
	keys := DiffToStringList2(d.Modified)
	sort.Sort(keys)
	for _, path := range keys {
		r.print("Modified path:", path)
		r.indent().printPath(d.Modified[path])
	}
}
func (r *Report) printPath(d *diff.PathDiff) {
	if !d.ExtensionsDiff.Empty() {
		r.printExtensionsDiff(d.ExtensionsDiff)
	}
	if !d.RefDiff.Empty() {
		r.printValue(d.RefDiff, "Reference")
	}
	if !d.SummaryDiff.Empty() {
		r.printValue(d.SummaryDiff, " Summary")
	}
	if !d.DescriptionDiff.Empty() {
		r.printValue(d.DescriptionDiff, "Description")
	}
	if !d.OperationsDiff.Empty() {
		r.printOperations(d.OperationsDiff)
	}
	if !d.ServersDiff.Empty() {
		r.printServers(d.ServersDiff)
	}
	if !d.ParametersDiff.Empty() {
		r.printParams(d.ParametersDiff)
	}
}
func (r *Report) printOperations(d *diff.OperationsDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Added)
	for _, added := range d.Added {
		r.print("New Operations Added:", added)
	}

	sort.Sort(d.Deleted)
	for _, deleted := range d.Deleted {
		r.print("Operations Deleted", deleted)
	}

	//problem with iterating
	keys := DiffToStringList3(d.Modified)
	sort.Sort(keys)
	for _, operation := range keys {
		r.print("Modified operation:", operation)
		r.indent().printMethod(d.Modified[operation])
	}
}
func (r *Report) printTagsDiff(d *diff.TagsDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Added)
	for _, added := range d.Added {
		r.print("New Tags Added:", added)
	}

	sort.Sort(d.Deleted)
	for _, deleted := range d.Deleted {
		r.print("Tags Deleted", deleted)
	}
	//problem with iterating
	keys := DiffToStringList4(d.Modified)
	sort.Sort(keys)
	for _, tags := range keys {
		r.print("Modified Tags:", tags)
		r.indent().printTags(d.Modified[tags])
	}
}
func (r *Report) printTagAdded(d *diff.TagsDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Added)
	for _, added := range d.Added {
		r.print("New Tags Added:", added)
	}
}
func (r *Report) printTagDeleted(d *diff.TagsDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Deleted)
	for _, deleted := range d.Deleted {
		r.print("Tags Deleted", deleted)
	}

}
func (r *Report) printTagModified(d *diff.TagsDiff) {
	if d.Empty() {
		return
	}
	keys := DiffToStringList4(d.Modified)
	sort.Sort(keys)
	for _, tags := range keys {
		r.print("Modified Tags:", tags)
		r.indent().printTags(d.Modified[tags])
	}
}
func (r *Report) printTags(d *diff.TagDiff) {
	if !d.NameDiff.Empty() {
		r.printValue(d.NameDiff, "Name")
	}
	if !d.DescriptionDiff.Empty() {
		r.printValue(d.DescriptionDiff, "Description")
	}
}
func (r *Report) printInfoDiff(d *diff.InfoDiff) {
	// ExtensionsDiff     *ExtensionsDiff
	// ContactDiff        *ContactDiff
	// LicenseDiff        *LicenseDiff
	// TitleDiff          *ValueDiff
	// DescriptionDiff    *ValueDiff
	// TermsOfServiceDiff *ValueDiff
	// VersionDiff        *ValueDiff
	if !d.ContactDiff.Empty() {
		if d.ContactDiff.Added {
			r.print("Contact Diff Added")
		}
		if d.ContactDiff.Deleted {
			r.print("Contact Diff Deleted")
		}
		if !d.ContactDiff.NameDiff.Empty() {
			r.printValue(d.ContactDiff.NameDiff, "Names")
		}
		if !d.ContactDiff.URLDiff.Empty() {
			r.printValue(d.ContactDiff.URLDiff, "URL")
		}
		if !d.ContactDiff.EmailDiff.Empty() {
			r.printValue(d.ContactDiff.EmailDiff, "Email")
		}
		if !d.ContactDiff.ExtensionsDiff.Empty() {
			r.printExtensionsDiff(d.ContactDiff.ExtensionsDiff)
		}
	}
	if !d.LicenseDiff.Empty() {
		if d.LicenseDiff.Added {
			r.print("License Diff Added")
		}
		if d.LicenseDiff.Deleted {
			r.print("License Diff Deleted")
		}
		if !d.LicenseDiff.NameDiff.Empty() {
			r.printValue(d.LicenseDiff.NameDiff, "Names")
		}
		if !d.LicenseDiff.URLDiff.Empty() {
			r.printValue(d.LicenseDiff.URLDiff, "URL")
		}
		if !d.LicenseDiff.ExtensionsDiff.Empty() {
			r.printExtensionsDiff(d.LicenseDiff.ExtensionsDiff)
		}
	}
	if !d.TitleDiff.Empty() {
		r.printValue(d.TitleDiff, "Title")
	}
	if !d.DescriptionDiff.Empty() {
		r.printValue(d.DescriptionDiff, "Description")
	}
	if !d.TermsOfServiceDiff.Empty() {
		r.printValue(d.TermsOfServiceDiff, "Terms Of Service")
	}
	if !d.VersionDiff.Empty() {
		r.printValue(d.VersionDiff, "Version")
	}
}
func (r *Report) printEndpointsAdded(d *diff.EndpointsDiff) {
	if d.Empty() {
		return
	}
	if len(d.Added) != 0 {
		r.printTitle("New Endpoints", len(d.Added))
		sort.Sort(d.Added)
		for _, added := range d.Added {
			r.print(added.Method, added.Path, " ")
		}
		r.print("")
	}
}
func (r *Report) PrintEndpointsDeleted(d *diff.Endpoint) {

	IsBreaking = true
	r.print(d.Method, d.Path, " ")
	r.print("")
}
func (r *Report) PrintEndpointsModified(key diff.Endpoint, value *diff.MethodDiff) {
	IsBreaking = true
	r.print(key.Method, key.Path)
	r.indent().printMethod(value)
	r.print("")
}
func (r *Report) printServerAdded(d *diff.ServersDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Added)
	for _, added := range d.Added {
		r.print("New server:", added)
	}

}
func (r *Report) printServerDeleted(d *diff.ServersDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Deleted)
	for _, deleted := range d.Deleted {
		r.print("Deleted server:", deleted)
	}
}
func (r *Report) printServerModified(d *diff.ServersDiff) {
	if d.Empty() {
		return
	}
	keys := diff.ToStringList(d.Modified)
	sort.Sort(keys)
	for _, server := range keys {
		r.print("Modified server:", server)
		r.indent().printServer(d.Modified[server])
	}
}
func (r *Report) printServers(d *diff.ServersDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Added)
	for _, added := range d.Added {
		r.print("New server:", added)
	}

	sort.Sort(d.Deleted)
	for _, deleted := range d.Deleted {
		r.print("Deleted server:", deleted)
	}

	keys := diff.ToStringList(d.Modified)
	sort.Sort(keys)
	for _, server := range keys {
		r.print("Modified server:", server)
		r.indent().printServer(d.Modified[server])
	}
}

func (r *Report) printMethod(d *diff.MethodDiff) {
	if d.Empty() {
		return
	}

	r.printValue(d.DescriptionDiff, "Description")
	r.printParams(d.ParametersDiff)

	if !d.RequestBodyDiff.Empty() {
		r.print("Request body changed")
		r.indent().printRequestBody(d.RequestBodyDiff)
	}

	if !d.ResponsesDiff.Empty() {
		r.print("Responses changed")
		r.indent().printResponses(d.ResponsesDiff)
	}

	r.printMessage(d.CallbacksDiff, "Callbacks changed")
	r.printValue(d.DeprecatedDiff, "Deprecated")

	if !d.SecurityDiff.Empty() {
		r.print("Security changed")
		r.indent().printSecurityRequirements(d.SecurityDiff)
	}

	if !d.ServersDiff.Empty() {
		r.print("Servers changed")
		r.indent().printServers(d.ServersDiff)
	}
}

func (r *Report) printParams(d *diff.ParametersDiff) {
	if d.Empty() {
		return
	}

	for _, location := range diff.ParamLocations {
		params := d.Added[location]
		sort.Strings(params)
		for _, param := range params {
			r.print("New", location, "param:", param)
		}
	}

	for _, location := range diff.ParamLocations {
		params := d.Deleted[location]
		sort.Strings(params)
		for _, param := range params {
			r.print("Deleted", location, "param:", param)
		}
	}

	for _, location := range diff.ParamLocations {
		paramDiffs := d.Modified[location]
		keys := diff.ToStringList(paramDiffs)
		sort.Sort(keys)
		for _, param := range keys {
			r.print("Modified", location, "param:", param)
			r.indent().printParam(paramDiffs[param])
		}
	}
}

func (r *Report) printParam(d *diff.ParameterDiff) {
	r.printValue(d.DescriptionDiff, "Description")
	r.printValue(d.StyleDiff, "Style")
	r.printValue(d.ExplodeDiff, "Explode")
	r.printValue(d.AllowEmptyValueDiff, "AllowEmptyValue")
	r.printValue(d.AllowReservedDiff, "AllowReserved")
	r.printValue(d.DeprecatedDiff, "Deprecated")
	r.printValue(d.RequiredDiff, "Required")

	if !d.SchemaDiff.Empty() {
		r.print("Schema changed")
		r.indent().printSchema(d.SchemaDiff)
	}

	r.printValue(d.ExampleDiff, "Example")

	if !d.ExamplesDiff.Empty() {
		r.print("Examples changed")
		r.indent().printExamples(d.ExamplesDiff)
	}

	if !d.ContentDiff.Empty() {
		r.print("Content changed")
		r.indent().printContent(d.ContentDiff)
	}
}

func (r *Report) printExamples(d *diff.ExamplesDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Added)
	for _, example := range d.Added {
		r.print("New example:", example)
	}

	sort.Sort(d.Deleted)
	for _, example := range d.Deleted {
		r.print("Deleted example:", example)
	}

	keys := diff.ToStringList(d.Modified)
	sort.Sort(keys)
	for _, example := range keys {
		r.print("Modified example:", example)
		r.indent().printExample(d.Modified[example])
	}
}

func (r *Report) printExample(d *diff.ExampleDiff) {
	if d.Empty() {
		return
	}

	r.printValue(d.SummaryDiff, "Summary")
	r.printValue(d.DescriptionDiff, "Description")
	r.printValue(d.ValueDiff, "Value")
	r.printValue(d.ExternalValueDiff, "ExternalValue")
}

func (r *Report) printRequiredProperties(d *diff.RequiredPropertiesDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Added)
	for _, added := range d.Added {
		r.print("New required property:", added)
	}

	sort.Sort(d.Deleted)
	for _, deleted := range d.Deleted {
		r.print("Deleted required property:", deleted)
	}
}

func (r *Report) printServer(d *diff.ServerDiff) {
	if d.Empty() {
		return
	}

	r.printConditional(d.Added, "Server added")
	r.printConditional(d.Deleted, "Server deleted")
	r.printValue(d.URLDiff, "URL")
	r.printValue(d.DescriptionDiff, "Description")
	if !d.VariablesDiff.Empty() {
		r.print("Variables changed")
		r.indent().printVariables(d.VariablesDiff)
	}
}

func (r *Report) printVariables(d *diff.VariablesDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Added)
	for _, variable := range d.Added {
		r.print("New variable:", variable)
	}

	sort.Sort(d.Deleted)
	for _, variable := range d.Deleted {
		r.print("Deleted variable:", variable)
	}

	keys := diff.ToStringList(d.Modified)
	sort.Sort(keys)
	for _, variable := range keys {
		r.print("Modified variable:", variable)
		r.indent().printVariable(d.Modified[variable])
	}
}

func (r *Report) printVariable(d *diff.VariableDiff) {
	if d.Empty() {
		return
	}

	if !d.EnumDiff.Empty() {
		r.printConditional(len(d.EnumDiff.Added) > 0, "New enum values:", d.EnumDiff.Added)
		r.printConditional(len(d.EnumDiff.Deleted) > 0, "Deleted enum values:", d.EnumDiff.Deleted)
	}
	r.printValue(d.DefaultDiff, "Default")
	r.printValue(d.DescriptionDiff, "Description")
}

func (r *Report) printSchema(d *diff.SchemaDiff) {
	if d.Empty() {
		return
	}

	r.printConditional(d.SchemaAdded, "Schema added")
	r.printConditional(d.SchemaDeleted, "Schema deleted")
	r.printConditional(d.CircularRefDiff, "Schema circular referecnce changed")

	if !d.OneOfDiff.Empty() {
		r.print("Property 'OneOf' changed")
		r.indent().printSchemaListDiff(d.OneOfDiff)
	}
	if !d.AnyOfDiff.Empty() {
		r.print("Property 'AnyOf' changed")
		r.indent().printSchemaListDiff(d.AnyOfDiff)
	}
	if !d.AllOfDiff.Empty() {
		r.print("Property 'AllOf' changed")
		r.indent().printSchemaListDiff(d.AllOfDiff)
	}

	if !d.NotDiff.Empty() {
		r.print("Property 'Not' changed")
		r.indent().printSchema(d.NotDiff)
	}

	r.printValue(d.TypeDiff, "Type")
	r.printValue(d.TitleDiff, "Title")
	r.printValue(d.FormatDiff, "Format")
	r.printValue(d.DescriptionDiff, "Description")

	if !d.EnumDiff.Empty() {
		r.printConditional(len(d.EnumDiff.Added) > 0, "New enum values:", d.EnumDiff.Added)
		r.printConditional(len(d.EnumDiff.Deleted) > 0, "Deleted enum values:", d.EnumDiff.Deleted)
	}

	r.printValue(d.DefaultDiff, "Default")
	r.printValue(d.ExampleDiff, "Example")
	r.printValue(d.AdditionalPropertiesAllowedDiff, "AdditionalProperties")
	r.printValue(d.UniqueItemsDiff, "UniqueItems")
	r.printValue(d.ExclusiveMinDiff, "ExclusiveMin")
	r.printValue(d.ExclusiveMaxDiff, "ExclusiveMax")
	r.printValue(d.NullableDiff, "Nullable")
	r.printValue(d.ReadOnlyDiff, "ReadOnly")
	r.printValue(d.WriteOnlyDiff, "WriteOnly")
	r.printValue(d.AllowEmptyValueDiff, "AllowEmptyValue")
	r.printValue(d.XMLDiff, "XML")
	r.printValue(d.DeprecatedDiff, "Deprecated")
	r.printValue(d.MinDiff, "Min")
	r.printValue(d.MaxDiff, "Max")
	r.printValue(d.MultipleOfDiff, "MultipleOf")
	r.printValue(d.MinLengthDiff, "MinLength")
	r.printValue(d.MaxLengthDiff, "MaxLength")
	r.printValue(d.PatternDiff, "Pattern")
	r.printValue(d.MinItemsDiff, "MinItems")
	r.printValue(d.MaxItemsDiff, "MaxItems")

	if !d.ItemsDiff.Empty() {
		r.print("Items changed")
		r.indent().printSchema(d.ItemsDiff)
	}

	if !d.RequiredDiff.Empty() {
		r.print("Required changed")
		r.indent().printRequiredProperties(d.RequiredDiff)
	}

	r.printValue(d.MinPropsDiff, "MinProps")
	r.printValue(d.MaxPropsDiff, "MaxProps")

	if !d.PropertiesDiff.Empty() {
		r.print("Properties changed")
		r.indent().printProperties(d.PropertiesDiff)
	}

	if !d.AdditionalPropertiesDiff.Empty() {
		r.print("AdditionalProperties changed")
		r.indent().printSchema(d.AdditionalPropertiesDiff)
	}

	r.printMessage(d.DiscriminatorDiff, "Discriminator changed")
}

func (r *Report) printSchemaListDiff(d *diff.SchemaListDiff) {
	if d.Empty() {
		return
	}

	if d.Added > 0 {
		r.print(d.Added, "schemas added")
	}
	if d.Deleted > 0 {
		r.print(d.Deleted, "schemas deleted")
	}
	if len(d.Modified) > 0 {
		for schemaRef, schemaDiff := range d.Modified {
			r.print("Schema", schemaRef, "modified")
			r.indent().printSchema(schemaDiff)
		}
	}
}

func (r *Report) printProperties(d *diff.SchemasDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Added)
	for _, property := range d.Added {
		r.print("New property:", property)
	}

	sort.Sort(d.Deleted)
	for _, property := range d.Deleted {
		r.print("Deleted property:", property)
	}

	keys := diff.ToStringList(d.Modified)
	sort.Sort(keys)
	for _, property := range keys {
		r.print("Modified property:", property)
		r.indent().printSchema(d.Modified[property])
	}
}

func quote(value interface{}) interface{} {
	if value == nil {
		return "null"
	}
	if reflect.ValueOf(value).Kind() == reflect.String {
		return "'" + value.(string) + "'"
	}
	return value
}

func (r *Report) printResponses(d *diff.ResponsesDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Added)
	for _, added := range d.Added {
		r.print("New response:", added)
	}

	sort.Sort(d.Deleted)
	for _, deleted := range d.Deleted {
		r.print("Deleted response:", deleted)
	}

	keys := diff.ToStringList(d.Modified)
	sort.Sort(keys)
	for _, response := range keys {
		r.print("Modified response:", response)
		r.indent().printResponse(d.Modified[response])
	}
}

func (r *Report) printResponse(d *diff.ResponseDiff) {
	if d.Empty() {
		return
	}

	r.printValue(d.DescriptionDiff, "Description")

	if !d.ContentDiff.Empty() {
		r.print("Content changed")
		r.indent().printContent(d.ContentDiff)
	}

	if !d.HeadersDiff.Empty() {
		r.print("Headers changed")
		r.indent().printHeaders(d.HeadersDiff)
	}
}

func (r *Report) printRequestBody(d *diff.RequestBodyDiff) {
	if d.Empty() {
		return
	}

	r.printValue(d.DescriptionDiff, "Description")

	if !d.ContentDiff.Empty() {
		r.print("Content changed")
		r.indent().printContent(d.ContentDiff)
	}
}

func (r *Report) printContent(d *diff.ContentDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.MediaTypeAdded)
	for _, name := range d.MediaTypeAdded {
		r.print("New media type:", name)
	}

	sort.Sort(d.MediaTypeDeleted)
	for _, name := range d.MediaTypeDeleted {
		r.print("Deleted media type:", name)
	}

	keys := diff.ToStringList(d.MediaTypeModified)
	sort.Sort(keys)
	for _, name := range keys {
		r.print("Modified media type:", name)
		r.indent().printMediaType(d.MediaTypeModified[name])
	}
}

func (r *Report) printMediaType(d *diff.MediaTypeDiff) {
	if d.Empty() {
		return
	}

	if !d.SchemaDiff.Empty() {
		r.print("Schema changed")
		r.indent().printSchema(d.SchemaDiff)
	}

	r.printValue(d.ExampleDiff, "Example")

	if !d.ExamplesDiff.Empty() {
		r.print("Examples changed")
		r.indent().printExamples(d.ExamplesDiff)
	}

	r.printMessage(d.EncodingsDiff, "Encodings changed")
}

func (r *Report) printValue(d *diff.ValueDiff, title string) {
	if d.Empty() {
		return
	}

	r.print(title, "changed from", quote(d.From), "to", quote(d.To))
}
func (r *Report) printHeaders(d *diff.HeadersDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Added)
	for _, added := range d.Added {
		r.print("New header:", added)
	}

	sort.Sort(d.Deleted)
	for _, deleted := range d.Deleted {
		r.print("Deleted header:", deleted)
	}

	keys := diff.ToStringList(d.Modified)
	sort.Sort(keys)
	for _, header := range keys {
		r.print("Modified header:", header)
		r.indent().printHeader(d.Modified[header])
	}
}

func (r *Report) printHeader(d *diff.HeaderDiff) {
	if d.Empty() {
		return
	}

	r.printValue(d.DescriptionDiff, "Description")
	r.printValue(d.DeprecatedDiff, "Deprecated")
	r.printValue(d.RequiredDiff, "Required")

	r.printValue(d.ExampleDiff, "Example")

	if !d.ExamplesDiff.Empty() {
		r.print("Examples changed")
		r.indent().printExamples(d.ExamplesDiff)
	}

	if !d.SchemaDiff.Empty() {
		r.print("Schema changed")
		r.indent().printSchema(d.SchemaDiff)
	}

	if !d.ContentDiff.Empty() {
		r.print("Content changed")
		r.indent().printContent(d.ContentDiff)
	}
}

func (r *Report) printSecurityAdded(d *diff.SecurityRequirementsDiff) {
	if d.Empty() {
		return
	}
	sort.Sort(d.Added)
	for _, added := range d.Added {
		r.print("New security requirements:", added)
		IsBreaking = true
	}
}
func (r *Report) printSecurityDeleted(d *diff.SecurityRequirementsDiff) {
	if d.Empty() {
		return
	}
	sort.Sort(d.Deleted)
	for _, deleted := range d.Deleted {
		r.print("Deleted security requirements:", deleted)
		IsBreaking = true
	}
}
func (r *Report) printSecurityModified(d *diff.SecurityRequirementsDiff) {
	if d.Empty() {
		return
	}
	keys := diff.ToStringList(d.Modified)
	sort.Sort(keys)
	for _, securityRequirementID := range keys {
		r.print("Modified security requirements:", securityRequirementID)
		r.indent().printSecurityScopes(d.Modified[securityRequirementID])
		IsBreaking = true
	}
}
func (r *Report) printSecurityRequirements(d *diff.SecurityRequirementsDiff) {
	if d.Empty() {
		return
	}

	sort.Sort(d.Added)
	for _, added := range d.Added {
		r.print("New security requirements:", added)
		IsBreaking = true
	}

	sort.Sort(d.Deleted)
	for _, deleted := range d.Deleted {
		r.print("Deleted security requirements:", deleted)
		IsBreaking = true
	}

	keys := diff.ToStringList(d.Modified)
	sort.Sort(keys)
	for _, securityRequirementID := range keys {
		r.print("Modified security requirements:", securityRequirementID)
		r.indent().printSecurityScopes(d.Modified[securityRequirementID])
		IsBreaking = true
	}
}

func (r *Report) printSecurityScopes(d diff.SecurityScopesDiff) {
	keys := diff.ToStringList(d)
	sort.Sort(keys)
	for _, scheme := range keys {
		scopeDiff := d[scheme]
		r.printConditional(len(scopeDiff.Added) > 0, "Scheme", scheme, "Added scopes:", scopeDiff.Added)
		r.printConditional(len(scopeDiff.Deleted) > 0, "Scheme", scheme, "Deleted scopes:", scopeDiff.Deleted)
	}
}

func (r *Report) printTitle(title string, count int) {
	text := ""
	if count == 0 {
		text = fmt.Sprintf("### %s: None", title)
	} else {
		text = fmt.Sprintf("### %s: %d", title, count)
	}

	r.print(text)
	r.print(strings.Repeat("-", len(text)))
}

func (r *Report) printMessage(d diff.IDiff, output ...interface{}) {
	r.printConditional(!d.Empty(), output...)
}

func (r *Report) printConditional(b bool, output ...interface{}) {
	if b {
		r.print(output...)
	}
}
func DiffToStringList1(m diff.ModifiedInterfaces) diff.StringList {
	keys := make(diff.StringList, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
func DiffToStringList2(m diff.ModifiedPaths) diff.StringList {
	keys := make(diff.StringList, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
func DiffToStringList3(m diff.ModifiedOperations) diff.StringList {
	keys := make(diff.StringList, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
func DiffToStringList4(m diff.ModifiedTags) diff.StringList {
	keys := make(diff.StringList, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
