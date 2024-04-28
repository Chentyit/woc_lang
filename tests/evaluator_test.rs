#[cfg(test)]
mod evaluator_test {
    use woc_lang::{
        evaluator::evaluator::eval,
        object::object::{BaseValue, Object, Value},
        parser_v2::parser::Parser,
    };

    #[test]
    fn test_return_in_block_stmt() {
        let tests = vec![
            (
                "{ return 10; }",
                Object::Return(BaseValue::Integer(Value::new(10))),
            ),
            (
                "{ return 10; return 20; }",
                Object::Return(BaseValue::Integer(Value::new(10))),
            ),
            (
                "{ return 10; 20; }",
                Object::Return(BaseValue::Integer(Value::new(10))),
            ),
            (
                "{ 10; return 20; }",
                Object::Return(BaseValue::Integer(Value::new(20))),
            ),
            (
                "{ 10; 20; return 30; }",
                Object::Return(BaseValue::Integer(Value::new(30))),
            ),
            (
                "if (1 > 2) { 10; } else { 20; }",
                Object::Base(BaseValue::Integer(Value::new(20))),
            ),
            (
                "if (1 > 2) { 10; } else { 20; return 30; }",
                Object::Return(BaseValue::Integer(Value::new(30))),
            ),
            (
                "if (10 > 1) { if (10 > 1) { return 30; } return 1; }",
                Object::Return(BaseValue::Integer(Value::new(30))),
            ),
        ];

        for (input, expected) in tests {
            let evaluated = test_eval(input);
            test_equal_object(evaluated, expected);
        }
    }

    #[test]
    fn test_return_stmt() {
        let tests = vec![
            (
                "return 1 && 1;",
                Object::Return(BaseValue::Boolean(Value::new(true))),
            ),
            (
                "return 1 || 1;",
                Object::Return(BaseValue::Boolean(Value::new(true))),
            ),
            (
                "return 10;",
                Object::Return(BaseValue::Integer(Value::new(10))),
            ),
            (
                "return 8.22;",
                Object::Return(BaseValue::Float(Value::new(8.22))),
            ),
            (
                "return true;",
                Object::Return(BaseValue::Boolean(Value::new(true))),
            ),
            (
                "return false;",
                Object::Return(BaseValue::Boolean(Value::new(false))),
            ),
            (
                "return 1 + 1;",
                Object::Return(BaseValue::Integer(Value::new(2))),
            ),
            (
                "return 1.1 + 1.1;",
                Object::Return(BaseValue::Float(Value::new(2.2))),
            ),
            (
                "return 1 == 1;",
                Object::Return(BaseValue::Boolean(Value::new(true))),
            ),
            (
                "return 1 != 1;",
                Object::Return(BaseValue::Boolean(Value::new(false))),
            ),
            (
                "return 1 < 1;",
                Object::Return(BaseValue::Boolean(Value::new(false))),
            ),
            (
                "return 1 > 1;",
                Object::Return(BaseValue::Boolean(Value::new(false))),
            ),
            (
                "return 1 <= 1;",
                Object::Return(BaseValue::Boolean(Value::new(true))),
            ),
            (
                "return 1 >= 1;",
                Object::Return(BaseValue::Boolean(Value::new(true))),
            ),
            (
                "return !1;",
                Object::Return(BaseValue::Boolean(Value::new(false))),
            ),
            (
                "return !!1;",
                Object::Return(BaseValue::Boolean(Value::new(true))),
            ),
            (
                "return !true;",
                Object::Return(BaseValue::Boolean(Value::new(false))),
            ),
            (
                "return !!true;",
                Object::Return(BaseValue::Boolean(Value::new(true))),
            ),
            (
                "return !false;",
                Object::Return(BaseValue::Boolean(Value::new(true))),
            ),
            (
                "return !!false;",
                Object::Return(BaseValue::Boolean(Value::new(false))),
            ),
        ];

        for (input, expected) in tests {
            let evaluated = test_eval(input);
            test_equal_object(evaluated, expected);
        }
    }

    #[test]
    fn test_if_exp() {
        let tests = vec![
            (
                "if (true) { 10; }",
                Object::Base(BaseValue::Integer(Value::new(10))),
            ),
            (
                "if (1) { 10; }",
                Object::Base(BaseValue::Integer(Value::new(10))),
            ),
            (
                "if (1 < 2) { 10; }",
                Object::Base(BaseValue::Integer(Value::new(10))),
            ),
            (
                "if (1 > 2) { 10; } else { 20; }",
                Object::Base(BaseValue::Integer(Value::new(20))),
            ),
            (
                "if (1 < 2) { 10; } else { 20; }",
                Object::Base(BaseValue::Integer(Value::new(10))),
            ),
            ("if (false) { 10; }", Object::Null),
            ("if (1 > 2) { 10; }", Object::Null),
            (
                "if (1 < 2) { 10; } else if ( 2 > 1 ) { 20; } else { 30; }",
                Object::Base(BaseValue::Integer(Value::new(10))),
            ),
            (
                "if (1 == 2) { 10; } else if ( 2 > 1 ) { 20; } else { 30; }",
                Object::Base(BaseValue::Integer(Value::new(20))),
            ),
            (
                "if (1 == 2) { 10; } else if ( 2 < 1 ) { 20; } else { 30; }",
                Object::Base(BaseValue::Integer(Value::new(30))),
            ),
        ];

        for (input, expected) in tests {
            let evaluated = test_eval(input);
            test_equal_object(evaluated, expected);
        }
    }

    #[test]
    fn test_eval_bool_exp() {
        let tests = vec![
            ("true;", true),
            ("false;", false),
            ("1 < 2;", true),
            ("1 > 2;", false),
            ("1 < 1;", false),
            ("1 > 1;", false),
            ("1 == 1;", true),
            ("1 != 1;", false),
            ("1 == 2;", false),
            ("1 != 2;", true),
            ("true == true;", true),
            ("false == false;", true),
            ("true == false;", false),
            ("true != false;", true),
            ("false != true;", true),
            ("(1 < 2) == true;", true),
            ("(1 < 2) == false;", false),
            ("(1 > 2) == true;", false),
            ("(1 > 2) == false;", true),
            ("true && true;", true),
            ("true && false;", false),
            ("false && true;", false),
            ("false && false;", false),
            ("true || true;", true),
            ("true || false;", true),
            ("false || true;", true),
            ("false || false;", false),
            ("(1 < 2) && true;", true),
            ("(1 < 2) && false;", false),
            ("(1 > 2) && true;", false),
            ("(1 > 2) && false;", false),
            ("(1 < 2) || true;", true),
            ("(1 < 2) || false;", true),
            ("(1 > 2) || true;", true),
            ("(1 > 2) || false;", false),
            ("(1 < 2) && (1 == 1);", true),
            ("(1 < 2) && (1 != 1);", false),
            ("(1 > 2) && (1 == 1);", false),
            ("(1 > 2) && (1 != 1);", false),
            ("(1 < 2) || (1 == 1);", true),
            ("(1 < 2) || (1 != 1);", true),
            ("(1 > 2) || (1 == 1);", true),
            ("(1 > 2) || (1 != 1);", false),
            ("!1;", false),
            ("!0;", true),
            ("!!1;", true),
            ("!!0;", false),
            ("1 && true", true),
            ("1 && false", false),
            ("0 && true", false),
            ("0 && false", false),
            ("1 || true", true),
            ("1 || false", true),
            ("0 || true", true),
            ("0 || false", false),
            ("!1;", false),
            ("!1 && true;", false),
            ("!1 || true;", true),
            ("!0;", true),
            ("!0 && true;", true),
            ("!0 || true;", true),
            ("!!1;", true),
            ("!!1 && true;", true),
            ("!!1 || true;", true),
            ("!!0;", false),
            ("!!0 && true;", false),
            ("!!0 || true;", true),
        ];

        for (input, expected) in tests {
            let evaluated = test_eval(input);
            test_equal_object(
                evaluated,
                Object::Base(BaseValue::Boolean(Value::new(expected))),
            );
        }
    }

    #[test]
    fn test_eval_prefix_exp() {
        let _tests = vec![
            ("!true;", false),
            ("!false;", true),
            ("!5;", false),
            ("!!true;", true),
            ("!!false;", false),
            ("!!5;", true),
        ];

        for (input, expected) in _tests {
            let evaluated = test_eval(input);
            test_equal_object(
                evaluated,
                Object::Base(BaseValue::Boolean(Value::new(expected))),
            );
        }
    }

    #[test]
    fn test_eval_integer_exp() {
        let _tests = vec![("5;", 5), ("10;", 10), ("-5;", -5), ("-10;", -10)];

        for (input, expected) in _tests {
            let value = test_eval(input);
            test_equal_object(
                value,
                Object::Base(BaseValue::Integer(Value::new(expected))),
            );
        }
    }

    #[test]
    fn test_eval_float_exp() {
        let _tests = vec![
            ("5.5;", 5.5),
            ("10.5;", 10.5),
            ("-5.5;", -5.5),
            ("-10.5;", -10.5),
            ("5 + 5.5;", 10.5),
            ("5.5 + 5;", 10.5),
            ("5.5 + 5.5;", 11.0),
            ("5 - 5.5;", -0.5),
            ("5.5 - 5;", 0.5),
            ("5.5 - 5.5;", 0.0),
            ("5 * 5.5;", 27.5),
            ("5.5 * 5;", 27.5),
            ("5.5 * 5.5;", 30.25),
            ("5.5 / 5;", 1.1),
            ("5.5 / 5.5;", 1.0),
        ];

        for (input, expected) in _tests {
            let evaluated = test_eval(input);
            test_equal_object(
                evaluated,
                Object::Base(BaseValue::Float(Value::new(expected))),
            );
        }
    }

    #[test]
    fn test_eval_boolean_exp() {
        let _tests = vec![("true;", true), ("false;", false)];

        for (input, expected) in _tests {
            let evaluated = test_eval(input);
            test_equal_object(
                evaluated,
                Object::Base(BaseValue::Boolean(Value::new(expected))),
            );
        }
    }

    fn test_eval(input: &str) -> Object {
        let parser = Parser::new(input);
        let program = parser.programs();

        return eval(program.get(0).unwrap());
    }

    fn test_equal_object(value: Object, expected: Object) {
        let get = value.clone();
        let want = expected.clone();

        if (get.is_null()) && (want.is_null()) {
            return;
        }

        match (value, expected) {
            (Object::Base(BaseValue::Integer(v)), Object::Base(BaseValue::Integer(e))) => {
                assert_eq!(v.value(), e.value());
            }
            (Object::Base(BaseValue::Float(v)), Object::Base(BaseValue::Float(e))) => {
                assert_eq!(v.value(), e.value());
            }
            (Object::Base(BaseValue::Boolean(v)), Object::Base(BaseValue::Boolean(e))) => {
                assert_eq!(v.value(), e.value());
            }
            (Object::Return(BaseValue::Integer(v)), Object::Return(BaseValue::Integer(e))) => {
                assert_eq!(v.value(), e.value());
            }
            (Object::Return(BaseValue::Float(v)), Object::Return(BaseValue::Float(e))) => {
                assert_eq!(v.value(), e.value());
            }
            (Object::Return(BaseValue::Boolean(v)), Object::Return(BaseValue::Boolean(e))) => {
                assert_eq!(v.value(), e.value());
            }
            _ => panic!("The object is not equal, got={:?}, want={:?}", get, want),
        }
    }
}
