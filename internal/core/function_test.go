package core

// func TestFindStructsReturnedByFunction(t *testing.T) {
// 	// sample/basic パッケージを読み込む
// 	cfg := &packages.Config{
// 		Mode: packages.NeedTypes |
// 			packages.NeedSyntax |
// 			packages.NeedTypesInfo |
// 			packages.NeedName |
// 			packages.NeedImports,
// 		Dir: "../../sample/basic/repository",
// 	}

// 	pkgs, err := packages.Load(cfg, ".")
// 	if err != nil {
// 		t.Fatalf("パッケージの読み込みに失敗しました: %v", err)
// 	}

// 	if len(pkgs) == 0 {
// 		t.Fatal("パッケージが見つかりませんでした")
// 	}

// 	if packages.PrintErrors(pkgs) > 0 {
// 		t.Fatal("パッケージにエラーがあります")
// 	}

// 	pkg := pkgs[0]

// 	t.Run("NewConfig関数は*Configを返す", func(t *testing.T) {
// 		// NewConfig 関数を取得
// 		obj := pkg.Types.Scope().Lookup("NewConfig")
// 		if obj == nil {
// 			t.Fatal("NewConfig関数が見つかりませんでした")
// 		}

// 		fn, ok := obj.(*types.Func)
// 		if !ok {
// 			t.Fatal("NewConfigが関数ではありません")
// 		}

// 		// return 文から実際に返される型を取得
// 		returnTypes := FindStructsReturnedByFunction(pkg, fn)

// 		if len(returnTypes) == 0 {
// 			t.Fatal("返り値の型が取得できませんでした")
// 		}

// 		// *Config 型が返されることを確認
// 		firstType := returnTypes[0]
// 		ptrType, ok := firstType.(*types.Pointer)
// 		if !ok {
// 			t.Fatalf("ポインタ型ではありません: %v", firstType)
// 		}

// 		namedType, ok := ptrType.Elem().(*types.Named)
// 		if !ok {
// 			t.Fatalf("Named型ではありません: %v", ptrType.Elem())
// 		}

// 		if namedType.Obj().Name() != "Config" {
// 			t.Errorf("型名が期待値と異なります。期待: Config, 実際: %s", namedType.Obj().Name())
// 		}
// 	})

// 	t.Run("NewUserRepository関数はUserRepositoryインターフェースを返す(シグネチャ上)", func(t *testing.T) {
// 		// NewUserRepository 関数を取得
// 		obj := pkg.Types.Scope().Lookup("NewUserRepository")
// 		if obj == nil {
// 			t.Fatal("NewUserRepository関数が見つかりませんでした")
// 		}

// 		fn, ok := obj.(*types.Func)
// 		if !ok {
// 			t.Fatal("NewUserRepositoryが関数ではありません")
// 		}

// 		// return 文から実際に返される型を取得
// 		returnTypes := FindStructsReturnedByFunction(pkg, fn)

// 		if len(returnTypes) == 0 {
// 			t.Fatal("返り値の型が取得できませんでした")
// 		}

// 		// 実際には *userRepositoryImpl が返される
// 		firstType := returnTypes[0]

// 		// ポインタ型を確認
// 		ptrType, ok := firstType.(*types.Pointer)
// 		if !ok {
// 			t.Logf("返り値の型: %v (%T)", firstType, firstType)
// 			// インターフェース型の場合もあるのでエラーにしない
// 			return
// 		}

// 		namedType, ok := ptrType.Elem().(*types.Named)
// 		if !ok {
// 			t.Fatalf("Named型ではありません: %v", ptrType.Elem())
// 		}

// 		if namedType.Obj().Name() != "userRepositoryImpl" {
// 			t.Errorf("型名が期待値と異なります。期待: userRepositoryImpl, 実際: %s", namedType.Obj().Name())
// 		}
// 	})

// 	t.Run("存在しない関数を渡した場合は空のスライスを返す", func(t *testing.T) {
// 		// ダミーの関数を作成
// 		dummyPkg := types.NewPackage("test/dummy", "dummy")
// 		sig := types.NewSignature(nil, nil, nil, false)
// 		dummyFn := types.NewFunc(0, dummyPkg, "DummyFunc", sig)

// 		returnTypes := FindStructsReturnedByFunction(pkg, dummyFn)

// 		if len(returnTypes) != 0 {
// 			t.Errorf("空のスライスが期待されますが、長さが %d です", len(returnTypes))
// 		}
// 	})
// }
