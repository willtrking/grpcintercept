package main

import (
  "strings"
  "bytes"
  "flag"
  "fmt"
  "io/ioutil"
  "go/format"
  "bufio"
  "regexp"
  "text/template"
  "os"
)


const _pkg_prefix = "package "
const _mult_imp_ = "import ("

var (
	service = flag.String("Service", "", "The receiver passed to grpc")
  grpcRegister = flag.String("GRPCRegister","","The function to register with grpc")
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type TemplateData struct {
  Package string
  SourceFile string
  ServiceType string
  RegisterFunc string
  IsMainFile bool
  ProtobufPackages []string
  InterceptorPackages []string
  ServiceCalls []ServiceDef
}

type ServiceDef struct {
  CallName string
  ServiceType string
  InputProto string
  OutputProto string
  InputProtoName string
  InputInterceptor string
}

func isServiceFunc(st string,tmpldat TemplateData) bool {

  rxp := fmt.Sprintf("^func \\([A-z_]+\\s(\\*%s|%s)\\)",tmpldat.ServiceType,tmpldat.ServiceType)
  re := regexp.MustCompile(rxp)
  return re.MatchString(st)
}

func stripFuncReceiver(st string,tmpldat TemplateData) string {
  rxp := fmt.Sprintf("^func \\([A-z_]+\\s(\\*%s|%s)\\)",tmpldat.ServiceType,tmpldat.ServiceType)
  re := regexp.MustCompile(rxp)
  return re.ReplaceAllString(st,"")
}

func stripComments(rpl []byte) []byte {
  re := regexp.MustCompile("(?s)//.*?\n|/\\*.*?\\*/")
  return re.ReplaceAll(rpl,nil)
}

func cleanupImport(st string) string {
  first := strings.TrimPrefix(st,"import ")

  return strings.TrimSpace(strings.Trim(first,"\""))
}

func cleanupFuncDef(st string) string {

  r := strings.NewReplacer(")","","(","","{","","}","")

  f := r.Replace(st)
  return strings.TrimSpace(f)
}

func stringInSlice(str string, list []string) bool {
 for _, v := range list {
   if v == str {
     return true
   }
 }
 return false
}

func funcArgSplit(st string) (nme string, tpe string) {

  split := strings.Split(strings.TrimSpace(st), " ")

  if len(split) == 1 {
    return "",split[0]
  }else if len(split) == 2{
    return split[0],split[1]
  }else{
    panic("Function arguments not 1/2")
  }

}

func typePackage(st string) string {
  split := strings.Split(strings.TrimSpace(st), ".")

  split[0] = strings.TrimPrefix(split[0],"*")
  split[0] = strings.TrimPrefix(split[0],"&")
  return split[0]
}

func formatInputPackage(name string, pkg string) string {
  return fmt.Sprintf("%s \"%s\"",name,pkg)
}

func main(){
  flag.Parse()

  info,_ := os.Stat(flag.Args()[0])
  inputFileMode := info.Mode()

  dat, err := ioutil.ReadFile(flag.Args()[0])
  check(err)
  formatted,ferr := format.Source(dat)
  check(ferr)


  formatted = stripComments(formatted)

  scanner := bufio.NewScanner(bytes.NewReader(formatted))

  tmpldat := TemplateData{
    SourceFile: flag.Args()[0],
    ServiceType: *service,
    RegisterFunc: *grpcRegister,
    IsMainFile: false,
  }

  insideImport := false
  importMap := map[string]string{}

  var curFuncName string
  insideFuncInput := false
  insideFuncOutput := false
  funcDefinitionMap := map[string]string{}


  /*Scan through our formatted ource
    Comments are already stripped out here
    Goal is to generate a map of service functions
    with their definitions, and a map of imports to their packages

    In addition we're figuring out the package name
    AND we're checking if this is the main file of the package
    To determine that, we look if a type with our ServiceType name
    is defined in the source. If it is, we're the main
  */
  for scanner.Scan() {
    //fmt.Println("----")
    line := scanner.Text()
    //fmt.Println(scanner.Text())

    if !insideFuncInput {
      if strings.HasPrefix(line,"type "+tmpldat.ServiceType){
        tmpldat.IsMainFile = true
        continue
      }

      if isServiceFunc(line,tmpldat){

        pure := stripFuncReceiver(line,tmpldat)
        curFuncName = strings.TrimSpace(pure[:strings.Index(pure,"(")])
        line = pure[strings.Index(pure,"(")+1:]
        funcDefinitionMap[curFuncName] = ""
        insideFuncInput = true

      }
    }

    addedDef := false
    if insideFuncInput {
      line = strings.TrimSpace(line)
      funcDefinitionMap[curFuncName] = fmt.Sprintf("%s %s",funcDefinitionMap[curFuncName],line)
      addedDef = true
      foundOutput := false
      for _,ir := range line {
        if ir == ')'{
          insideFuncInput = false
        }

        if !insideFuncInput{
          if ir == '(' {
            foundOutput = true
          }
        }
      }

      if !insideFuncInput && !foundOutput {
        delete(funcDefinitionMap,curFuncName)
        curFuncName = ""
      }else if !insideFuncInput && foundOutput {
        insideFuncOutput = true
        line = line[strings.Index(line,"(")+1:]
      }

    }

    if insideFuncOutput {
      line = strings.TrimSpace(line)
      if !addedDef{
        funcDefinitionMap[curFuncName] = fmt.Sprintf("%s %s",funcDefinitionMap[curFuncName],line)
      }
      hasBodyDef := false
      for _,or := range line {
        if or == ')'{
          insideFuncOutput = false
        }

        if !insideFuncOutput {
          if or == '{' {
            hasBodyDef = true
          }
        }

      }

      if !insideFuncOutput && !hasBodyDef {
        delete(funcDefinitionMap,curFuncName)
        curFuncName = ""
      }else if !insideFuncOutput {
        curFuncName = ""
      }

      continue
    }




    //Trim line and do checks
    line = strings.TrimSpace(scanner.Text())

    if line == ""{
      continue
    }

    if !insideImport{
      if strings.HasPrefix(line, _pkg_prefix) {
        tmpldat.Package = strings.TrimPrefix(line,_pkg_prefix)
        continue
      }

      if strings.HasPrefix(line,_mult_imp_) {
        insideImport = true
      }else if strings.HasPrefix(line,"import "){
        pkg := cleanupImport(line)
        split := strings.Split(pkg,"/")
        rename := split[len(split)-1]
        importMap[rename] = pkg
      }


    }else{
      if strings.HasSuffix(line,")"){
        insideImport = false
        continue
      }

      //fmt.Println(line)
      if strings.Contains(line," \""){
        //Renaming import
        split := strings.SplitAfterN(line," \"",2)
        rename := cleanupImport(split[0])
        pkg := cleanupImport(split[1])
        importMap[rename] = pkg

      }else{
        pkg := cleanupImport(line)
        split := strings.Split(pkg,"/")
        rename := split[len(split)-1]
        importMap[rename] = pkg


      }
    }


	}

  /* Figure out which function definitions are really relevant
    gRPC/protobuf has a specific format for functions, so we'll follow along with that

    In protobuf, an RPC service HAS to take in 1 arg and return 1 arg

    In grpc, the stubs generated from the proto files will have 2 arguments
    in and 2 arguments out. The second argument returned HAS to be an error,
    and the first argument taken in HAS to be context.
    The remaining arguments are the service input/response.
    In addition, this func HAS to have a receiver of our ServiceType type

    To check if this should be intercepted, we need to see if the function
    has a total of 3 input arguments, and matches all of the above.

  */


  for funcName, funcDef := range funcDefinitionMap {
    split := strings.SplitAfterN(funcDef,") (",2)
    input := cleanupFuncDef(split[0])
    output := cleanupFuncDef(split[1])

    splitInput := strings.Split(input,",")
    splitOutput := strings.Split(output,",")
    if len(splitOutput) != 2{
      //gRPC requires 2 returns
      continue
    }

    if len(splitInput) < 2 {
      //Require at least 2 inputs
      continue
    }

    _,lstTpe := funcArgSplit(splitOutput[1])
    if lstTpe != "error" {
      //Last argument must be error
      continue
    }

    _,frstTpe := funcArgSplit(splitInput[0])

    if !strings.Contains(frstTpe,".Context"){
      //First argument is not context
      continue
    }

    //Yes, this needs to be done here
    if len(splitInput) != 3 {
      panic("Found service definition, but not enough input args for middleware")
    }

    inputProtoNme,inputProtoTpe := funcArgSplit(splitInput[1])
    inputPbPkg := typePackage(inputProtoTpe)
    inputPbPkgImp := formatInputPackage(inputPbPkg,importMap[inputPbPkg])

    _,outputProtoTpe := funcArgSplit(splitOutput[0])
    outputPbPkg := typePackage(outputProtoTpe)
    outputPbPkgImp := formatInputPackage(outputPbPkg,importMap[outputPbPkg])

    if !stringInSlice(inputPbPkgImp,tmpldat.ProtobufPackages){
        tmpldat.ProtobufPackages = append(tmpldat.ProtobufPackages,inputPbPkgImp)
    }
    if !stringInSlice(outputPbPkgImp,tmpldat.ProtobufPackages){
        tmpldat.ProtobufPackages = append(tmpldat.ProtobufPackages,outputPbPkgImp)
    }

    _,inputIntTpe := funcArgSplit(splitInput[2])
    inputIntPkg := typePackage(inputIntTpe)
    inputIntPkgImp := formatInputPackage(inputIntPkg,importMap[inputIntPkg])

    if !stringInSlice(inputIntPkgImp,tmpldat.InterceptorPackages){
        tmpldat.InterceptorPackages = append(tmpldat.InterceptorPackages,inputIntPkgImp)
    }


    serviceDef := new(ServiceDef)
    serviceDef.ServiceType = tmpldat.ServiceType
    serviceDef.CallName = funcName
    serviceDef.InputProto=strings.TrimSpace(splitInput[1])
    serviceDef.InputProtoName=inputProtoNme
    serviceDef.OutputProto=strings.TrimSpace(splitOutput[0])
    serviceDef.InputInterceptor=strings.TrimSpace(inputIntTpe)

    tmpldat.ServiceCalls = append(tmpldat.ServiceCalls,*serviceDef)

  }

  registerPkg := typePackage(tmpldat.RegisterFunc)
  registerPkgImp := formatInputPackage(registerPkg,importMap[registerPkg])

  if !stringInSlice(registerPkgImp,tmpldat.ProtobufPackages){
      tmpldat.ProtobufPackages = append(tmpldat.ProtobufPackages,registerPkgImp)
  }

  var toWrite bytes.Buffer

  tmpl := template.Must(template.New("outputProtoTpe").Parse(TEMPLATE))
  tmpl.Execute(&toWrite,tmpldat)

  final,finalErr := format.Source(toWrite.Bytes())
  check(finalErr)


  outputFile := tmpldat.SourceFile[:strings.LastIndex(tmpldat.SourceFile,".")]
  outputFile = outputFile + ".intr.go"

  writeErr := ioutil.WriteFile(outputFile, final, inputFileMode)
  check(writeErr)
}
