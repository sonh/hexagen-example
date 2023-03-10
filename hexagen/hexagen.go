package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	trimprefix  = flag.String("trimprefix", "", "trim the `prefix` from the generated constant names")
	linecomment = flag.Bool("linecomment", false, "use line comment text as printed text when present")
	buildTags   = flag.String("tags", "", "comma-separated list of build tags to apply")

	// typeNames = flag.String("type", "", "comma-separated list of type names; must be set")
	// output = flag.String("output", "", "output file name; default srcdir/<type>_string.go")
	// outpkg = flag.String("outpkg", "", "output file name; default srcdir/<type>_string.go")

	typeNames string
	output    string
	outpkg    string
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of usrgen:\n")
	fmt.Fprintf(os.Stderr, "\tusrgen [flags] -type T [directory]\n")
	fmt.Fprintf(os.Stderr, "\tusrgen [flags] -type T files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

var rootCmd = &cobra.Command{
	Use:   "hexagen",
	Short: "hexagen is a code generator for abstractions of hexagon architecture",
	Long:  `hexagen is a code generator for abstractions of hexagon architecture`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("afsdf")
	},
}

var entMockCmd = &cobra.Command{
	Use:   "ent-mock",
	Short: "generate mock entity from entity abstraction",
	Long:  `generate mock entity from entity abstraction`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hexagen generating ent mock for %s...\n", typeNames)
		genMock(typeNames, output, outpkg, args...)
		fmt.Printf("hexagen genterated ent mock for %s\n", typeNames)
	},
}

var entImplCmd = &cobra.Command{
	Use:   "ent-impl",
	Short: "generate entity impl from entity abstraction",
	Long:  `generate entity impl from entity abstraction`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hexagen generating ent impl for %s...\n", typeNames)
		genEntImpl(typeNames, output, outpkg, args...)
		fmt.Printf("hexagen genterated ent impl for %s\n", typeNames)
	},
}

var entRepoCmd = &cobra.Command{
	Use:   "ent-repo",
	Short: "generate entity impl from entity abstraction",
	Long:  `generate entity impl from entity abstraction`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hexagen generating ent repo for %s...\n", typeNames)
		genEntRepo(typeNames, output, outpkg, args...)
		fmt.Printf("hexagen genterated ent repo for %s\n", typeNames)
	},
}

var entProtoCmd = &cobra.Command{
	Use:   "ent-proto",
	Short: "generate entity impl from entity abstraction",
	Long:  `generate entity impl from entity abstraction`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hexagen generating impl for %s\n", typeNames)
		genEntProto(typeNames, output, outpkg, args...)
		fmt.Printf("hexagen genterated impl for %s\n", typeNames)
	},
}

func main() {
	entMockCmd.Flags().StringVarP(&typeNames, "type", "t", "", "")
	entMockCmd.Flags().StringVarP(&output, "output", "o", ".", "")
	entMockCmd.Flags().StringVarP(&outpkg, "outpkg", "p", "", "")
	if err := entMockCmd.MarkFlagRequired("type"); err != nil {
		panic(err)
	}

	entImplCmd.Flags().StringVarP(&typeNames, "type", "t", "", "")
	entImplCmd.Flags().StringVarP(&output, "output", "o", ".", "")
	entImplCmd.Flags().StringVarP(&outpkg, "outpkg", "p", "", "")
	if err := entImplCmd.MarkFlagRequired("type"); err != nil {
		panic(err)
	}

	entProtoCmd.Flags().StringVarP(&typeNames, "type", "t", "", "")
	entProtoCmd.Flags().StringVarP(&output, "output", "o", ".", "")
	entProtoCmd.Flags().StringVarP(&outpkg, "outpkg", "p", "", "")
	if err := entProtoCmd.MarkFlagRequired("type"); err != nil {
		panic(err)
	}

	entRepoCmd.Flags().StringVarP(&typeNames, "type", "t", "", "")
	entRepoCmd.Flags().StringVarP(&output, "output", "o", ".", "")
	entRepoCmd.Flags().StringVarP(&outpkg, "outpkg", "p", "", "")
	if err := entRepoCmd.MarkFlagRequired("type"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(entMockCmd)
	rootCmd.AddCommand(entImplCmd)
	rootCmd.AddCommand(entProtoCmd)
	rootCmd.AddCommand(entRepoCmd)
	rootCmd.Execute()
	return
}

func genMock(typeNamesText string, output string, outpkg string, args ...string) {
	/*if len(typeNamesText) == 0 {
		flag.Usage()
		os.Exit(2)
	}*/

	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	// We accept either one directory or a list of files. Which do we have?
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	var outputDir string
	if len(args) == 1 && isDirectory(args[0]) {
		outputDir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		outputDir = filepath.Dir(args[0])
	}

	g := Generator{
		trimPrefix:  *trimprefix,
		lineComment: *linecomment,
	}
	g.parsePackage(args, tags)

	if outpkg == "" {
		outpkg = g.pkg.name
	}

	// Print the header and package clause.
	g.Printf("// Code generated by \"hexagen %s %s\"; DO NOT EDIT.\n", typeNamesText, strings.Join(args, " "))
	g.Printf("\n")
	g.Printf("package %s", outpkg)
	g.Printf("\n")

	importPaths := []string{
		"\"backend/internal/usermgmt/modules/user/core/entity\"\n",
		"\"backend/internal/usermgmt/pkg/field\"\n",
		"\"fmt\"\n",
	}
	if len(importPaths) == 1 {
		g.Printf("import %s", importPaths[0])
	} else {
		g.Printf("import (\n")
		for _, importPath := range importPaths {
			g.Printf(importPath)
		}
		g.Printf(")\n")
	}

	// Run generateMock for each type.
	typeNames := strings.Split(typeNamesText, ",")
	for _, typeName := range typeNames {
		g.generateMock(typeName)
	}

	// Format the output.
	src := g.format()

	// Write to file.
	baseName := fmt.Sprintf("mock_%s.go", typeNames[0])
	if output == "" {
		output = outputDir
	}
	dest := filepath.Join(output, strings.ToLower(baseName))

	err := os.WriteFile(dest, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

// Value represents a declared constant.
type Value struct {
	originalName string // The name of the constant.
	name         string // The name with trimmed prefix.
	// The value is stored as a bit pattern alone. The boolean tells us
	// whether to interpret it as an int64 or a uint64; the only place
	// this matters is when sorting.
	// Much of the time the str field is all we need; it is printed
	// by Value.String.
	value  uint64 // Will be converted to int64 when needed.
	signed bool   // Whether the constant is a signed type.
	str    string // The string representation given by the "go/constant" package.
}

func (v *Value) String() string {
	return v.str
}

type Package struct {
	name  string
	defs  map[*ast.Ident]types.Object
	files []*File
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// These fields are reset for each type being generated.
	typeName string  // Name of the constant type.
	values   []Value // Accumulator for constant values of that type.

	trimPrefix  string
	lineComment bool

	InterfaceDecls []InterfaceDecl
}

type InterfaceDecl struct {
	Name string

	Methods []InterfaceMethodDecl
}

type InterfaceMethodDecl struct {
	Name       string
	ResultType []string
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	genDecl, ok := node.(*ast.GenDecl)
	if !ok || genDecl.Tok != token.TYPE {
		// We only care about type declarations.
		// ex:
		// 	type User interface {}
		// 	type User struct {}
		return true
	}

	/*// The name of the type of the constants we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by values.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).*/

	interfaceDecls := make([]InterfaceDecl, 0, len(genDecl.Specs))

	for _, spec := range genDecl.Specs {
		typeSpec := spec.(*ast.TypeSpec) // Guaranteed to succeed as this is TypeSpec.

		if typeSpec.Name.Name != f.typeName {
			// This is not the type we're looking for.
			continue
		}

		if typeSpec.Type != nil {
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			interfaceDecl := InterfaceDecl{
				Name:    typeSpec.Name.Name,
				Methods: make([]InterfaceMethodDecl, 0, len(interfaceType.Methods.List)),
			}

			// Interface's methods
			for _, method := range interfaceType.Methods.List {
				interfaceMethodDecl := InterfaceMethodDecl{
					Name:       method.Names[0].Name,
					ResultType: make([]string, 0, 3),
				}

				methodType, ok := method.Type.(*ast.FuncType)
				if !ok {
					continue
				}

				for _, methodResult := range methodType.Results.List {
					methodSelectorExpr, ok := methodResult.Type.(*ast.SelectorExpr)
					if !ok {
						continue
					}
					methodSelectorExprX, ok := methodSelectorExpr.X.(*ast.Ident)
					resultType := fmt.Sprintf("%s.%s", methodSelectorExprX, methodSelectorExpr.Sel.Name)

					interfaceMethodDecl.ResultType = append(interfaceMethodDecl.ResultType, resultType)
				}

				interfaceDecl.Methods = append(interfaceDecl.Methods, interfaceMethodDecl)
			}
			interfaceDecls = append(interfaceDecls, interfaceDecl)
		}
	}

	f.InterfaceDecls = interfaceDecls
	return false
}

type Generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *Package     // Package we are scanning.

	trimPrefix  string
	lineComment bool
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func (g *Generator) parsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	g.addPackage(pkgs[0])
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file:        file,
			pkg:         g.pkg,
			trimPrefix:  g.trimPrefix,
			lineComment: g.lineComment,
		}
	}
}

// generateMock produces the String method for the named type.
func (g *Generator) generateMock(typeName string) {
	interfaceDecls := make([]InterfaceDecl, 0, 100)

	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			interfaceDecls = append(interfaceDecls, file.InterfaceDecls...)
		}
	}

	if len(interfaceDecls) == 0 {
		log.Fatalf(`no interface declarations defined for type "%s"`, typeName)
	}

	for _, interfaceDecl := range interfaceDecls {
		entityName := fmt.Sprintf("Valid%s", interfaceDecl.Name)

		g.Printf("type %s struct {\n", entityName)
		g.Printf("\tRandomID string\n")
		g.Printf("}\n")
		g.Printf("//This statement will fail to compile if *%s ever stops matching the interface.\n", interfaceDecl.Name)
		g.Printf("var _ entity.User = (*%s)(nil)\n", entityName)
		for _, interfaceMethod := range interfaceDecl.Methods {
			receiver := strings.ToLower(interfaceDecl.Name)
			receiverType := entityName
			funcResult := strings.Join(interfaceMethod.ResultType, ", ")
			if len(interfaceMethod.ResultType) > 1 {
				funcResult = fmt.Sprintf("(%s)", funcResult)
			}

			switch interfaceMethod.Name {
			case "UserID":
				g.Printf("func (%s %s) %s() %s {\n", receiver, receiverType, interfaceMethod.Name, funcResult)
				g.Printf("\treturn field.NewString(%s.RandomID)", receiver)
				g.Printf("}\n")
			case "Email":
				g.Printf("func (%s %s) %s() %s {\n", receiver, receiverType, interfaceMethod.Name, funcResult)
				g.Printf("\treturn field.NewString(fmt.Sprintf(\"%s+%s@example.com\", %s.RandomID))", LowerCaseFirstLetter(interfaceMethod.Name), "%s", receiver)
				g.Printf("}\n")
			default:
				g.Printf("func (%s %s) %s() %s {\n", receiver, receiverType, interfaceMethod.Name, funcResult)
				g.Printf("\treturn field.NewRandomString(64)")
				g.Printf("}\n")
			}
		}
	}

	for _, interfaceDecl := range interfaceDecls {
		validEntityName := fmt.Sprintf("Valid%s", interfaceDecl.Name)

		for _, interfaceMethod := range interfaceDecl.Methods {
			entityName := fmt.Sprintf("%sHasEmpty%s", interfaceDecl.Name, interfaceMethod.Name)
			receiver := strings.ToLower(interfaceDecl.Name)
			receiverType := entityName
			funcResult := strings.Join(interfaceMethod.ResultType, ", ")

			g.Printf("//%s keeps all valid attributes but override with an invalid %s\n", entityName, interfaceMethod.Name)
			g.Printf("type %s struct {\n", entityName)
			g.Printf("\t%s\n", validEntityName)
			g.Printf("}\n")
			g.Printf("var _ entity.User = (*%s)(nil)\n", entityName)

			g.Printf("func (%s %s) %s() %s {\n", receiver, receiverType, interfaceMethod.Name, funcResult)
			g.Printf("\treturn field.NewString(\"\")")
			g.Printf("}\n")
		}
	}
}

// ///////////////////////

func genEntImpl(typeNamesText string, output string, outpkg string, args ...string) {
	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	// We accept either one directory or a list of files. Which do we have?
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	var outputDir string
	if len(args) == 1 && isDirectory(args[0]) {
		outputDir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		outputDir = filepath.Dir(args[0])
	}

	g := Generator{
		trimPrefix:  *trimprefix,
		lineComment: *linecomment,
	}
	g.parsePackage(args, tags)

	if outpkg == "" {
		outpkg = g.pkg.name
	}

	// Print the header and package clause.
	g.Printf("// Code generated by \"hexagen %s %s\"; DO NOT EDIT.\n", typeNamesText, strings.Join(args, " "))
	g.Printf("\n")
	g.Printf("package %s", outpkg)
	g.Printf("\n")

	importPaths := []string{
		"\"github.com/pkg/errors\"\n",
		"\"backend/internal/usermgmt/pkg/field\"\n",
	}
	if len(importPaths) == 1 {
		g.Printf("import %s", importPaths[0])
	} else {
		g.Printf("import (\n")
		for _, importPath := range importPaths {
			g.Printf(importPath)
		}
		g.Printf(")\n")
	}

	// Run generateMock for each type.
	typeNames := strings.Split(typeNamesText, ",")
	for _, typeName := range typeNames {
		g.generateEntImpl(typeName)
	}

	// Format the output.
	src := g.format()

	// Write to file.
	baseName := fmt.Sprintf("%s_impl.go", typeNames[0])
	if output == "" {
		output = outputDir
	}
	dest := filepath.Join(output, strings.ToLower(baseName))

	err := os.WriteFile(dest, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

func (g *Generator) generateEntImpl(typeName string) {

	interfaceDecls := make([]InterfaceDecl, 0, 100)

	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			interfaceDecls = append(interfaceDecls, file.InterfaceDecls...)
		}
	}

	if len(interfaceDecls) == 0 {
		log.Fatalf(`no interface declarations defined for type "%s"`, typeName)
	}

	for _, interfaceDecl := range interfaceDecls {
		entityName := fmt.Sprintf("Null%s", interfaceDecl.Name)

		g.Printf("//This statement will fail to compile if *%s ever stops matching the interface.\n", entityName)
		g.Printf("var _ User = (*%s)(nil)\n", entityName)
		g.Printf("type %s struct {}\n", entityName)
		for _, interfaceMethod := range interfaceDecl.Methods {
			funcResult := strings.Join(interfaceMethod.ResultType, ", ")
			if len(interfaceMethod.ResultType) > 1 {
				funcResult = fmt.Sprintf("(%s)", funcResult)
			}
			g.Printf("func (%s %s) %s() %s {\n", strings.ToLower(interfaceDecl.Name), entityName, interfaceMethod.Name, funcResult)
			g.Printf("\treturn field.NewNullString()")
			g.Printf("}\n")
		}

	}

	for _, interfaceDecl := range interfaceDecls {
		funcName := fmt.Sprintf("Comapare%sValues", interfaceDecl.Name)
		argName1 := LowerCaseFirstLetter(interfaceDecl.Name) + "1"
		argName2 := LowerCaseFirstLetter(interfaceDecl.Name) + "2"

		g.Printf("\n")
		g.Printf("//%s compare values of two entities\n", funcName)
		g.Printf("func %s(%s %s, %s %s) error {\n", funcName, argName1, interfaceDecl.Name, argName2, interfaceDecl.Name)
		g.Printf("\tswitch {\n")
		for _, interfaceMethod := range interfaceDecl.Methods {
			g.Printf("\tcase %s.%s().%s() != %s.%s().%s():\n", argName1, interfaceMethod.Name, strings.Split(interfaceMethod.ResultType[0], ".")[1], argName2, interfaceMethod.Name, strings.Split(interfaceMethod.ResultType[0], ".")[1])
			g.Printf("\treturn errors.New(\"%s is not equal\")\n", interfaceMethod.Name)
		}
		g.Printf("\t}\n")
		g.Printf("\treturn nil\n")
		g.Printf("}\n")
	}

	for _, interfaceDecl := range interfaceDecls {
		sliceName := fmt.Sprintf("%ss", interfaceDecl.Name)

		g.Printf("// %s represents for a slice of %s\n", sliceName, interfaceDecl.Name)
		g.Printf("type %s []%s\n", sliceName, interfaceDecl.Name)
		/*g.Printf("}\n")
		g.Printf("//This statement will fail to compile if *%s ever stops matching the interface.\n", interfaceDecl.Name)
		g.Printf("var _ entity.User = (*User)(nil)\n")*/

		for _, interfaceMethod := range interfaceDecl.Methods {
			methodName := fmt.Sprintf("%ss", interfaceMethod.Name)
			methodResult := strings.Join(interfaceMethod.ResultType, ", ")
			receiver := strings.ToLower(sliceName)
			if len(interfaceMethod.ResultType) > 1 {
				methodResult = fmt.Sprintf("(%s)", methodResult)
			}
			g.Printf("func (%s %s) %s() []%s {\n", receiver, sliceName, methodName, methodResult)
			g.Printf("\t%s := make([]%s, 0, len(%s))\n", LowerCaseFirstLetter(methodName), methodResult, receiver)
			g.Printf("\tfor _, %s := range %s {\n", strings.ToLower(interfaceDecl.Name), receiver)
			g.Printf("\t\t%s = append(%s, %s.%s())", LowerCaseFirstLetter(methodName), LowerCaseFirstLetter(methodName), strings.ToLower(interfaceDecl.Name), interfaceMethod.Name)
			g.Printf("\t}\n")
			g.Printf("\treturn %s", LowerCaseFirstLetter(methodName))
			g.Printf("}\n")
			g.Printf("\n")
		}

	}
}

// /////////////////////////////////////////////////////////
func genEntRepo(typeNamesText string, output string, outpkg string, args ...string) {
	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	// We accept either one directory or a list of files. Which do we have?
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	var outputDir string
	if len(args) == 1 && isDirectory(args[0]) {
		outputDir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		outputDir = filepath.Dir(args[0])
	}

	g := Generator{
		trimPrefix:  *trimprefix,
		lineComment: *linecomment,
	}
	g.parsePackage(args, tags)

	if outpkg == "" {
		outpkg = "postgres" // g.pkg.name
	}

	// Print the header and package clause.
	g.Printf("// Code generated by \"hexagen %s %s\"; DO NOT EDIT.\n", typeNamesText, strings.Join(args, " "))
	g.Printf("\n")
	g.Printf("package %s", outpkg)
	g.Printf("\n")

	importPaths := []string{
		"\"context\"\n",
		"\"strings\"\n",
		"\"fmt\"\n",
		"\"backend/internal/usermgmt/pkg/field\"\n",
		"\"backend/internal/usermgmt/pkg/database\"\n",
		"\"backend/internal/usermgmt/modules/user/core/entity\"\n",
	}
	if len(importPaths) == 1 {
		g.Printf("import %s", importPaths[0])
	} else {
		g.Printf("import (\n")
		for _, importPath := range importPaths {
			g.Printf(importPath)
		}
		g.Printf(")\n")
	}

	// Run generateMock for each type.
	typeNames := strings.Split(typeNamesText, ",")
	for _, typeName := range typeNames {
		g.generateEntRepo(typeName)
	}

	// Format the output.
	src := g.format()

	// Write to file.
	baseName := fmt.Sprintf("%s_get_by.go", typeNames[0])
	if output == "" {
		output = outputDir
	}
	dest := filepath.Join(output, strings.ToLower(baseName))

	err := os.WriteFile(dest, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

func (g *Generator) generateEntRepo(typeName string) {

	interfaceDecls := make([]InterfaceDecl, 0, 100)

	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			interfaceDecls = append(interfaceDecls, file.InterfaceDecls...)
		}
	}

	if len(interfaceDecls) == 0 {
		log.Fatalf(`no interface declarations defined for type "%s"`, typeName)
	}

	for _, interfaceDecl := range interfaceDecls {
		/*sliceName := fmt.Sprintf("%ss", interfaceDecl.Name)
		g.Printf("type %s []%s\n", sliceName, interfaceDecl.Name)
		g.Printf("}\n")
		g.Printf("//This statement will fail to compile if *%s ever stops matching the interface.\n", interfaceDecl.Name)
		g.Printf("var _ entity.User = (*User)(nil)\n")*/

		for _, interfaceMethod := range interfaceDecl.Methods {
			methodName := "GetBy" + interfaceMethod.Name
			methodArg := LowerCaseFirstLetter(interfaceMethod.Name) + " " + interfaceMethod.ResultType[0]
			methodResult := fmt.Sprintf("(entity.%s, error)", interfaceDecl.Name)
			receiver := LowerCaseFirstLetter(interfaceDecl.Name + "Repo")
			receiverType := interfaceDecl.Name + "Repo"

			if len(interfaceMethod.ResultType) > 1 {
				methodResult = fmt.Sprintf("(%s)", methodResult)
			}
			g.Printf("func (%s *%s) %s(ctx context.Context, db database.QueryExecer, %s) %s {\n", receiver, receiverType, methodName, methodArg, methodResult)
			g.Printf("\t%s := &%s{}\n", strings.ToLower(interfaceDecl.Name), interfaceDecl.Name)
			g.Printf("fields := database.GetFieldNames(%s)\n", strings.ToLower(interfaceDecl.Name))
			g.Printf("\n")
			g.Printf("\tstmt := `SELECT %s FROM %s WHERE %s = $1`\n", "%s", "%s", "%s")
			g.Printf("\tstmt = fmt.Sprintf(stmt, strings.Join(fields, \",\"), %sTable%sColumn, %s.TableName())\n", interfaceDecl.Name, interfaceMethod.Name, strings.ToLower(interfaceDecl.Name))
			g.Printf("\n")
			g.Printf("\trow := db.QueryRow(ctx, stmt, %s)\n", LowerCaseFirstLetter(interfaceMethod.Name))
			g.Printf("\n")
			g.Printf("\tif err := row.Scan(database.GetScanFields(%s, fields)...); err != nil {\n", strings.ToLower(interfaceDecl.Name))
			g.Printf("\t\treturn nil, err")
			g.Printf("}\n")
			g.Printf("\n")
			g.Printf("\treturn %s, nil", strings.ToLower(interfaceDecl.Name))
			g.Printf("}\n")
			g.Printf("\n")
		}

	}
}

func LowerCaseFirstLetter(text string) string {
	if len(text) < 1 {
		return text
	}
	if len(text) == 1 {
		return strings.ToLower(text)
	}
	return strings.ToLower(string(text[0])) + text[1:]
}

// //////////////////////////////////////////
func genEntProto(typeNamesText string, output string, outpkg string, args ...string) {
	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	// We accept either one directory or a list of files. Which do we have?
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	var outputDir string
	if len(args) == 1 && isDirectory(args[0]) {
		outputDir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		outputDir = filepath.Dir(args[0])
	}

	g := Generator{
		trimPrefix:  *trimprefix,
		lineComment: *linecomment,
	}
	g.parsePackage(args, tags)

	if outpkg == "" {
		outpkg = g.pkg.name
	}

	// Print the header and package clause.
	g.Printf("// Code generated by \"hexagen %s %s\"; DO NOT EDIT.\n", typeNamesText, strings.Join(args, " "))
	g.Printf("\n")
	g.Printf("package %s", outpkg)
	g.Printf("\n")

	importPaths := []string{
		"\"backend/internal/usermgmt/pkg/field\"\n",
	}
	if len(importPaths) == 1 {
		g.Printf("import %s", importPaths[0])
	} else {
		g.Printf("import (\n")
		for _, importPath := range importPaths {
			g.Printf(importPath)
		}
		g.Printf(")\n")
	}

	// Run generateMock for each type.
	typeNames := strings.Split(typeNamesText, ",")
	for _, typeName := range typeNames {
		g.generateEntProto(typeName)
	}

	// Format the output.
	src := g.format()

	// Write to file.
	baseName := fmt.Sprintf("%s_ext.pb.go", typeNames[0])
	if output == "" {
		output = outputDir
	}
	dest := filepath.Join(output, strings.ToLower(baseName))

	err := os.WriteFile(dest, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

func (g *Generator) generateEntProto(typeName string) {

	interfaceDecls := make([]InterfaceDecl, 0, 100)

	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			interfaceDecls = append(interfaceDecls, file.InterfaceDecls...)
		}
	}

	if len(interfaceDecls) == 0 {
		log.Fatalf(`no interface declarations defined for type "%s"`, typeName)
	}

	for _, interfaceDecl := range interfaceDecls {
		/*sliceName := fmt.Sprintf("%ss", interfaceDecl.Name)
		g.Printf("type %s []%s\n", sliceName, interfaceDecl.Name)
		g.Printf("}\n")
		g.Printf("//This statement will fail to compile if *%s ever stops matching the interface.\n", interfaceDecl.Name)
		g.Printf("var _ entity.User = (*User)(nil)\n")*/

		for _, interfaceMethod := range interfaceDecl.Methods {
			methodName := interfaceMethod.Name
			methodResult := strings.Join(interfaceMethod.ResultType, ", ")
			receiver := LowerCaseFirstLetter(interfaceDecl.Name)
			receiverType := interfaceDecl.Name

			if len(interfaceMethod.ResultType) > 1 {
				methodResult = fmt.Sprintf("(%s)", methodResult)
			}
			g.Printf("func (%s *%s) %s() %s {\n", receiver, receiverType, methodName, methodResult)
			g.Printf("\treturn field.NewString(%s.%s)", receiver, methodName)
			g.Printf("}\n")
			g.Printf("\n")
		}

	}
}
