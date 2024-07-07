// Package logic provides shortcuts to create factors/CPTs for logic operations.
//
// To get logic factors, use either [Get] to get them by name,
// or one of the constructors of [Factor].
//
// Factor names equal constructor names,
// but are all lower case and use minus (-) instead of camel case.
// Examples:
//   - "And" becomes "and"
//   - "IfThen" becomes "if-then"
//
// Most logic factors work with binary True/False outcomes.
// The first possible outcome (index 0) is considered True,
// while the second outcome (index 1) is considered False.
package logic
