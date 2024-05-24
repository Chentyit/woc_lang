#[cfg(test)]
mod lexer_test {
    use woc_lang::{
        lexer::lexer::Lexer,
        token::token::Token,
    };
    use woc_lang::token::types::TokenType;

    #[test]
    fn test_demo() {
        let input = "let str = \"hello world\"; str;";

        let lexer = Lexer::new(input);
        let tokens_iter = lexer.tokens_iter();

        let expects = vec![
            Token::new(TokenType::Let, "let"),
            Token::new(TokenType::Ident, "str"),
            Token::new(TokenType::Assignment, "="),
            Token::new(TokenType::String, "hello world"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::Ident, "str"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::Eof, ""),
        ];

        for expect in expects {
            let token = tokens_iter.next().unwrap();
            assert_eq!(*expect.token_type(), *token.token_type());
            assert_eq!(expect.literal(), token.literal());
        }
    }

    #[test]
    fn test_string_token() {
        let input = "
            \"hello\"

            \"hello world\"

            \"hello
            world\"

            \"hello\\world\"

            \"hello 666 world\"
        ";

        let lexer = Lexer::new(input);
        let tokens_iter = lexer.tokens_iter();

        let expects = vec![
            Token::new(TokenType::String, "hello"),
            Token::new(TokenType::String, "hello world"),
            Token::new(TokenType::String, "hello\n            world"),
            Token::new(TokenType::String, "hello\\world"),
            Token::new(TokenType::String, "hello 666 world"),
            Token::new(TokenType::Eof, ""),
        ];

        for expect in expects {
            let token = tokens_iter.next().unwrap();
            assert_eq!(*expect.token_type(), *token.token_type());
            assert_eq!(expect.literal(), token.literal());
        }
    }

    #[test]
    fn test_ten_token() {
        let input = " ten  &&& ";

        let lexer = Lexer::new(input);
        let tokens_iter = lexer.tokens_iter();

        let tokens = vec![
            Token::new(TokenType::Ident, "ten"),
            Token::new(TokenType::And, "&&"),
            Token::new(TokenType::BitAnd, "&"),
            Token::new(TokenType::Eof, ""),
        ];

        for tok in tokens {
            let token = tokens_iter.next().unwrap();
            assert_eq!(*tok.token_type(), *token.token_type());
            assert_eq!(tok.literal(), token.literal());
        }
    }

    #[test]
    fn test_integer_and_float_num() {
        let input: &str = "10 1_000 1.0 1_000.0 99.99";

        let lexer = Lexer::new(input);
        let tokens_iter = lexer.tokens_iter();

        let tokens = vec![
            Token::new(TokenType::IntegerNum, "10"),
            Token::new(TokenType::IntegerNum, "1_000"),
            Token::new(TokenType::FloatNum, "1.0"),
            Token::new(TokenType::FloatNum, "1_000.0"),
            Token::new(TokenType::FloatNum, "99.99"),
            Token::new(TokenType::Eof, ""),
        ];

        for tok in tokens {
            let token = tokens_iter.next().unwrap();
            assert_eq!(*tok.token_type(), *token.token_type());
            assert_eq!(tok.literal(), token.literal());
        }
    }

    #[test]
    fn test_minus_number() {
        let input: &str = "-10 -1_000 -1.0 -1_000.0 -99.99";

        let lexer = Lexer::new(input);
        let tokens_iter = lexer.tokens_iter();

        let tokens = vec![
            Token::new(TokenType::IntegerNum, "-10"),
            Token::new(TokenType::IntegerNum, "-1_000"),
            Token::new(TokenType::FloatNum, "-1.0"),
            Token::new(TokenType::FloatNum, "-1_000.0"),
            Token::new(TokenType::FloatNum, "-99.99"),
            Token::new(TokenType::Eof, ""),
        ];

        for tok in tokens {
            let token = tokens_iter.next().unwrap();
            assert_eq!(*tok.token_type(), *token.token_type());
            assert_eq!(tok.literal(), token.literal());
        }
    }

    #[test]
    fn test_if_token() {
        let input = "
            if (5 < 10) {
                return true;
            ";

        let lexer = Lexer::new(input);
        let tokens_iter = lexer.tokens_iter();

        let tokens = vec![
            Token::new(TokenType::If, "if"),
            Token::new(TokenType::LeftParen, "("),
            Token::new(TokenType::IntegerNum, "5"),
            Token::new(TokenType::Less, "<"),
            Token::new(TokenType::IntegerNum, "10"),
            Token::new(TokenType::RightParen, ")"),
            Token::new(TokenType::LeftBrace, "{"),
            Token::new(TokenType::Return, "return"),
            Token::new(TokenType::True, "true"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::Eof, ""),
        ];

        for tok in tokens {
            let token = tokens_iter.next().unwrap();
            assert_eq!(*tok.token_type(), *token.token_type());
            assert_eq!(tok.literal(), token.literal());
        }
    }

    #[test]
    fn test_next_token() {
        let input = "
            let five = 5;
            let ten = 10;
            func (x, y) {
                x + y;
            }
            let result = add(five, ten);
            !-/*5;
            5 < 10 > 5;
            if (5 < 10) {
                return true;
            } else {
                return false;
            }
            10 == 10;
            10 != 9;
            for (let i = 0; i < 10; i += 1) {
                i;
            }
        ";

        let lexer = Lexer::new(input);
        let tokens_iter = lexer.tokens_iter();

        let tokens = vec![
            Token::new(TokenType::Let, "let"),
            Token::new(TokenType::Ident, "five"),
            Token::new(TokenType::Assignment, "="),
            Token::new(TokenType::IntegerNum, "5"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::Let, "let"),
            Token::new(TokenType::Ident, "ten"),
            Token::new(TokenType::Assignment, "="),
            Token::new(TokenType::IntegerNum, "10"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::Func, "func"),
            Token::new(TokenType::LeftParen, "("),
            Token::new(TokenType::Ident, "x"),
            Token::new(TokenType::Comma, ","),
            Token::new(TokenType::Ident, "y"),
            Token::new(TokenType::RightParen, ")"),
            Token::new(TokenType::LeftBrace, "{"),
            Token::new(TokenType::Ident, "x"),
            Token::new(TokenType::Plus, "+"),
            Token::new(TokenType::Ident, "y"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::RightBrace, "}"),
            Token::new(TokenType::Let, "let"),
            Token::new(TokenType::Ident, "result"),
            Token::new(TokenType::Assignment, "="),
            Token::new(TokenType::Ident, "add"),
            Token::new(TokenType::LeftParen, "("),
            Token::new(TokenType::Ident, "five"),
            Token::new(TokenType::Comma, ","),
            Token::new(TokenType::Ident, "ten"),
            Token::new(TokenType::RightParen, ")"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::Not, "!"),
            Token::new(TokenType::Minus, "-"),
            Token::new(TokenType::Slash, "/"),
            Token::new(TokenType::Asterisk, "*"),
            Token::new(TokenType::IntegerNum, "5"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::IntegerNum, "5"),
            Token::new(TokenType::Less, "<"),
            Token::new(TokenType::IntegerNum, "10"),
            Token::new(TokenType::Greater, ">"),
            Token::new(TokenType::IntegerNum, "5"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::If, "if"),
            Token::new(TokenType::LeftParen, "("),
            Token::new(TokenType::IntegerNum, "5"),
            Token::new(TokenType::Less, "<"),
            Token::new(TokenType::IntegerNum, "10"),
            Token::new(TokenType::RightParen, ")"),
            Token::new(TokenType::LeftBrace, "{"),
            Token::new(TokenType::Return, "return"),
            Token::new(TokenType::True, "true"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::RightBrace, "}"),
            Token::new(TokenType::Else, "else"),
            Token::new(TokenType::LeftBrace, "{"),
            Token::new(TokenType::Return, "return"),
            Token::new(TokenType::False, "false"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::RightBrace, "}"),
            Token::new(TokenType::IntegerNum, "10"),
            Token::new(TokenType::EqualTo, "=="),
            Token::new(TokenType::IntegerNum, "10"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::IntegerNum, "10"),
            Token::new(TokenType::NotEqualTo, "!="),
            Token::new(TokenType::IntegerNum, "9"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::For, "for"),
            Token::new(TokenType::LeftParen, "("),
            Token::new(TokenType::Let, "let"),
            Token::new(TokenType::Ident, "i"),
            Token::new(TokenType::Assignment, "="),
            Token::new(TokenType::IntegerNum, "0"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::Ident, "i"),
            Token::new(TokenType::Less, "<"),
            Token::new(TokenType::IntegerNum, "10"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::Ident, "i"),
            Token::new(TokenType::PlusAssign, "+="),
            Token::new(TokenType::IntegerNum, "1"),
            Token::new(TokenType::RightParen, ")"),
            Token::new(TokenType::LeftBrace, "{"),
            Token::new(TokenType::Ident, "i"),
            Token::new(TokenType::Semicolon, ";"),
            Token::new(TokenType::RightBrace, "}"),
            Token::new(TokenType::Eof, ""),
        ];

        for tok in tokens {
            let token = tokens_iter.next().unwrap();
            assert_eq!(*tok.token_type(), *token.token_type());
            assert_eq!(tok.literal(), token.literal());
        }
    }
}
