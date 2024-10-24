package esmongo

// For comparison of different BSON type values, see the [specified BSON comparison order].
//
// [specified BSON comparison order]: https://www.mongodb.com/docs/manual/reference/bson-type-comparison-order/#std-label-bson-types-comparison-order
type Operator string

const (
	// Comparison

	OperatorEQ  Operator = "$eq"  // Matches values that are equal to a specified value.
	OperatorGT  Operator = "$gt"  // Matches values that are greater than a specified value.
	OperatorGTE Operator = "$gte" // Matches values that are greater than or equal to a specified value.
	OperatorIN  Operator = "$in"  // Matches any of the values specified in an array.
	OperatorLT  Operator = "$lt"  // Matches values that are less than a specified value.
	OperatorLTE Operator = "$lte" // Matches values that are less than or equal to a specified value.
	OperatorNE  Operator = "$ne"  // Matches all values that are not equal to a specified value.
	OperatorNIN Operator = "$nin" // Matches none of the values specified in an array.

	// Logical

	OperatorAnd Operator = "$and" // Joins query clauses with a logical AND returns all documents that match the conditions of both clauses.
	OperatorNot Operator = "$not" // Inverts the effect of a query predicate and returns documents that do not match the query predicate.
	OperatorNor Operator = "$nor" // Joins query clauses with a logical NOR returns all documents that fail to match both clauses.
	OperatorOr  Operator = "$or"  // Joins query clauses with a logical OR returns all documents that match the conditions of either clause.

	// Element

	OperatorExists Operator = "$exists" // Matches documents that have the specified field.
	OperatorType   Operator = "$type"   // Selects documents if a field is of the specified type.

	// Evaluation

	OperatorExpr       Operator = "$expr"       // Allows use of aggregation expressions within the query language.
	OperatorJSONSchema Operator = "$jsonSchema" // Validate documents against the given JSON Schema.
	OperatorMod        Operator = "$mod"        // Performs a modulo operation on the value of a field and selects documents with a specified result.
	OperatorRegex      Operator = "$regex"      // Selects documents where values match a specified regular expression.
	OperatorText       Operator = "$text"       // Performs text search.
	OperatorWhere      Operator = "$where"      // Matches documents that satisfy a JavaScript expression.

	// Geospatial

	OperatorGeoInersects Operator = "$geoIntersects" // Selects geometries that intersect with a GeoJSON geometry. The 2dsphere index supports $geoIntersects.
	OperatorGeoWithin    Operator = "$geoWithin"     // Selects geometries within a bounding GeoJSON geometry. The 2dsphere and 2d indexes support $geoWithin.
	OperatorNear         Operator = "$near"          // Returns geospatial objects in proximity to a point. Requires a geospatial index. The 2dsphere and 2d indexes support $near.
	OperatorNearSphere   Operator = "$nearSphere"    // Returns geospatial objects in proximity to a point on a sphere. Requires a geospatial index. The 2dsphere and 2d indexes support $nearSphere.

	// Array

	OperatorAll       Operator = "$all"       // Matches arrays that contain all elements specified in the query.
	OperatorElemMatch Operator = "$elemMatch" // Selects documents if element in the array field matches all the specified $elemMatch conditions.
	OperatorSize      Operator = "$size"      // Selects documents if the array field is a specified size.

	// Bitwise

	OperatorBitsAlClear  Operator = "$bitsAllClear" // Matches numeric or binary values in which a set of bit positions all have a value of 0.
	OperatorBitsAllSet   Operator = "$bitsAllSet"   // Matches numeric or binary values in which a set of bit positions all have a value of 1.
	OperatorBitsAnyClear Operator = "$bitsAnyClear" // Matches numeric or binary values in which any bit from a set of bit positions has a value of 0.
	OperatorBitsAnySet   Operator = "$bitsAnySet"   // Matches numeric or binary values in which any bit from a set of bit positions has a value of 1.

	// Projection Operators

	OperatorFirst Operator = "$" // Projects the first element in an array that matches the query condition.
	// OperatorElemMatch Operator = "$elemMatch" // Projects the first element in an array that matches the specified $elemMatch condition.
	OperatorMeta  Operator = "$meta"  // Projects the document's score assigned during the $text operation.
	OperatorSlice Operator = "$slice" // Limits the number of elements projected from an array. Supports skip and limit slices.

	// Miscellaneous Operators

	OperatorRand    Operator = "$rand"    // Generates a random float between 0 and 1.
	OperatorNatural Operator = "$natural" // A special hint that can be provided via the sort() or hint() methods that can be used to force either a forward or reverse collection scan.
)
