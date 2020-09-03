module.exports = {
    "extends": [
        "standard"
    ],
    "plugins": [],
    "env": {
        "es6": true,
        "node": true
    },
    "parserOptions": {
        "sourceType": "module"
    },
    "globals": {
        "module": true,
        "Promise": true
    },
    "rules": {
        "curly": "error",
        "indent": [
            "warn",
            4,
            { "SwitchCase": 1 }
        ],
        "new-cap": [
            "error",
            {
                "capIsNew": false
            }
        ],
        "space-before-function-paren": "off",
        "semi": [
            "error",
            "always"
        ],
        "no-use-before-define": [
            "error",
            {
                "functions": false
            }
        ],
        "no-var": "error",
        "no-new": "error",
        "no-undef": "error",
        "no-unused-vars": "error",
        "prefer-const": "error",
        "no-const-assign": "error"
    },
    "overrides": [
        {
            "files": [
                "**/*.spec.js",
                "**/*.test.js"
            ],
            "env": {
                "jest": true
            },
            "globals": {
                "expect": true,
                "jest": true
            }
        }
    ]
}