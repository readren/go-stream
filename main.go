package main

import (
	"fmt"
	"os"

	typeCtors "github.com/readren/functional-go/typectors"
)

func main() {

	if false {
		s := typeCtors.BuildTypeName("map[string][3]image.Point")
		fmt.Println(s)

	} else {

		workingDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		config := typeCtors.Config{
			GeneratedPackageParentDir: fmt.Sprintf("%s/generated-package-parent-dir", workingDir),
			GeneratedPackageName:      "functional",
			TemplatesFolder:           fmt.Sprintf("%s/typectors", workingDir),
			TypesDescriptors: []typeCtors.TypeDescriptor{
				{
					TypeConstructorName: "stream",
					BaseTypeArguments:   typeCtors.TypeArguments{{Type: "int"}},
					TypeArgumentsForWhichPolymorphicMethodsAreInstantiated: []typeCtors.TypeArguments{
						{{Type: "image.Point", PackagePath: "image"}},
						{{Type: "bool"}, {Type: "int"}},
						{{Type: "string"}, {Type: "[]bool"}},
					},
				},
			},
		}

		typeCtors.GeneratePackage(config)
	}
}
