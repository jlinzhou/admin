module.exports = {
    root: true,
    env: {
        node: true
    },
    'extends': [
        'plugin:vue/essential',
        // '@vue/standard'
    ],
    rules: {
        'spaced-comment': 'off',
        'no-tabs': 'off',
        'indent': 'off',
        'space-before-function-paren': 'off',
        'no-trailing-spaces': 'off',
        'eol-last': 'off',
        'no-useless-return': 'off',
        'no-multiple-empty-lines': 'off',
        'prefer-promise-reject-errors': 'off',
        // 'import/no-duplicates': 'off',
        'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
        'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off'
    },
    parserOptions: {
        parser: 'babel-eslint'
    }
}
