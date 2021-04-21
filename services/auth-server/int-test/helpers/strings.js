module.exports = { compareStrings };

function compareStrings(a, b) {
    if (a < b) {
        return -1;
    }
    if (a > b) {
        return 1;
    }

    // names must be equal
    return 0;
}
