package demo

import (
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"testing"
)

func Test_demo1(t *testing.T) {
	//1.初始化一个新的评估环境，是一个包含变量(可以是自定义类型)、函数和类型信息的上下文，用于解析和评估表达式。环境声明变量和函数可以在表达式中被调用
	env, err := cel.NewEnv(
		cel.Variable("name", cel.StringType)) //参数类型绑定
	if err != nil {
		t.Fatal(err)
	}

	var str = `"Hello world! I'm " + name + "."`
	//2.用于解析和类型检查给定的表达式,并返回抽象语法树AST(包含env.Parse和env.Check)
	ast, issues := env.Compile(str) //编译、校验、执行str
	if issues.Err() != nil {
		t.Fatal(issues.Err())
	}

	//3.基于抽象语法树创建一个可执行的程序，将解析后的AST转换成一个可执行的对象，使得表达式可以在给定的上下文中被评估
	//编译AST(将类型检查后的AST编译成一个可执行程序) ->准备评估环境(配置和初始化评估环境，使得表达式能够在运行时进进行评估) ->返回程序对象(程序对象包含执行表达式所需的方法)
	program, err := env.Program(ast)
	if err != nil {
		t.Fatal(err)
	}

	//初始化name变量值
	values := map[string]interface{}{"name": "CEL"}
	//在指定的上下文中评估编译后的表达式
	eval, details, err := program.Eval(values)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(eval)
	fmt.Println(details)
}

func Test_demo2(t *testing.T) {

	env, err := cel.NewEnv(
		cel.Variable("i", cel.StringType),
		cel.Variable("you", cel.StringType),
		cel.Function("greet", cel.MemberOverload("string_greet_string", []*cel.Type{cel.StringType, cel.StringType}, cel.StringType,
			cel.BinaryBinding(func(lhs ref.Val, rhs ref.Val) ref.Val {
				return types.String(fmt.Sprintf("Hello %s! Nice to meet you, I'm %s.\n", rhs, lhs))
			}),
		),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	ast, issues := env.Compile("i.greet(you)")
	if issues.Err() != nil {
		t.Fatal(issues.Err())
	}

	program, err := env.Program(ast)
	if err != nil {
		t.Fatal(err)
	}

	values := map[string]interface{}{"i": "CEL",
		"you": "word"}

	eval, details, err := program.Eval(values)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(eval)
	fmt.Println(details)

}
