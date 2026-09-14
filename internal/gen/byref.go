package main

import (
	"fmt"
	"strings"
)

// A `get` whose resource can be fetched by reference takes `<ref>` rather
// than `<resource-id>`: an id, a CRN or a name, read by syntax exactly as the
// platform reads it. The SDK's Get<Resource>ByReference does the work; this
// file only shapes the command around it.
//
// Parents stay positional and stay ids, in the order the path takes them —
// `lb listener get <lb-id> <ref>` — because that is where the path already
// puts them. A name is unique only within its parent, so a resource whose
// list carries a parent filter (a subnet's --vpc) exposes that filter as a
// flag: without it a bare name can match more than one.

// emitByReferenceGet writes the `get <ref>` command for op.
func emitByReferenceGet(b *strings.Builder, svc service, r resourceGroup, op operation) error {
	ref := op.ByReference
	fn := fmt.Sprintf("new%s%s%sCommand", exported(svc.Name), exported(r.Name), exported(op.Verb))

	use := op.Verb
	for _, p := range op.PathParams[:len(op.PathParams)-1] {
		use += " <" + flagName(p.Wire) + ">"
	}
	use += " <ref>"

	short := op.Summary
	if short == "" {
		short = title(op.Verb) + " " + r.Name
	}
	short = strings.TrimSuffix(short, ".")

	kinds := "its id, its CRN or its name"
	if !ref.HasName {
		kinds = "its id or its CRN"
	}
	long := []string{
		short + ".",
		fmt.Sprintf("<ref> is %s. It is read by its syntax alone, the way the\n"+
			"platform reads it: a crn: value is a CRN, the 36-character UUID form is an\n"+
			"id, anything else is a name. A miss under one reading is not retried\n"+
			"under another.", kinds),
	}
	if !ref.HasName {
		long = append(long, "This resource has no name.")
	}
	scopeNames := make([]string, 0, len(ref.ScopeParams))
	for _, p := range ref.ScopeParams {
		name, err := resolveFlagName(op, p.Wire)
		if err != nil {
			return err
		}
		scopeNames = append(scopeNames, "--"+name)
	}
	if len(scopeNames) > 0 {
		long = append(long, fmt.Sprintf("A name is unique only within its parent: pass %s with a name, or the\n"+
			"lookup can match more than one and is refused.", strings.Join(scopeNames, " or ")))
	}

	fmt.Fprintf(b, "// %s builds `basaltic %s %s %s`.\n", fn, svc.Name, r.Name, op.Verb)
	fmt.Fprintf(b, "func %s(state *cli.State) *cobra.Command {\n", fn)
	fmt.Fprintf(b, "\tvar scope %s.%s\n", svc.Package, ref.ScopeType)
	fmt.Fprintf(b, "\tcmd := &cobra.Command{\n")
	fmt.Fprintf(b, "\t\tUse:   %q,\n", use)
	fmt.Fprintf(b, "\t\tShort: %q,\n", short)
	fmt.Fprintf(b, "\t\tArgs:  cobra.ExactArgs(%d),\n", len(op.PathParams))
	fmt.Fprintf(b, "\t\tLong: %s,\n", quote(strings.Join(long, "\n\n")))
	b.WriteString("\t\tRunE: func(cmd *cobra.Command, args []string) error {\n")
	fmt.Fprintf(b, "\t\t\tc, err := %sClient(state)\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n", unexported(svc.Name))

	callArgs := []string{"cmd.Context()"}
	for i := range op.PathParams {
		callArgs = append(callArgs, fmt.Sprintf("args[%d]", i))
	}
	callArgs = append(callArgs, "&scope")
	fmt.Fprintf(b, "\t\t\tout, err := c.%s(%s)\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n", ref.GoName, strings.Join(callArgs, ", "))
	b.WriteString("\t\t\treturn state.Printer().Value(out)\n")
	b.WriteString("\t\t},\n\t}\n")

	b.WriteString("\tf := cmd.Flags()\n\t_ = f\n")
	for _, p := range ref.ScopeParams {
		name, _ := resolveFlagName(op, p.Wire)
		fmt.Fprintf(b, "\tf.StringVar(&scope.%s, %q, \"\", %s)\n", p.GoName, name, quote(flagUsage(p)))
	}
	b.WriteString("\treturn cmd\n}\n\n")
	return nil
}
