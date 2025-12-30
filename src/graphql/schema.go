package graphql

// Schema defines the GraphQL schema
const Schema = `
type Query {
	health: String!
	speedTests: [SpeedTest!]!
	speedTest(id: ID!): SpeedTest
}

type SpeedTest {
	id: ID!
	timestamp: String!
	downloadMbps: Float!
	uploadMbps: Float!
	pingMs: Float!
	jitterMs: Float!
	packetLoss: Float!
	userAgent: String!
	shareCode: String
	shareViews: Int!
	createdAt: String!
}

type Mutation {
	startSpeedTest: SpeedTestStart!
}

type SpeedTestStart {
	testId: ID!
	status: String!
}
`
