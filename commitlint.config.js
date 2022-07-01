module.exports = {
    extends: ['@commitlint/config-conventional'],
    ignores: [
        (message) => message.includes('Initial commit')
    ],
    'type-enum': ['feat', 'fix', 'ci', 'docs', 'lint', 'test'],
};