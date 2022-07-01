module.exports = {
    extends: ['@commitlint/config-conventional'],
    ignores: [
        (message) => message.includes('Initial commit')
    ],
    rules: {
        'type-enum': [
            2,
            'always',
            [
            'feat', 'fix', 'docs', 'lint', 'ci', 'test', 'revert',
            ],
        ],
    },
};