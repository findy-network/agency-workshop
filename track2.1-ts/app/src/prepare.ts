import { agencyv1, AgentClient } from '@findy-network/findy-common-ts'
import { existsSync, readFileSync, writeFileSync } from 'fs';

export default async (agentClient: AgentClient, tag: string) => {

  const waitForSchema = async (schemaId: string) => new Promise<void>((resolve) => {
    const getSchema = async () => {
      const schemaMsg = new agencyv1.Schema();
      schemaMsg.setId(schemaId)
      try {
        await agentClient.getSchema(schemaMsg)
        resolve()
      } catch {
        setTimeout(getSchema, 1000, schemaId)
        return
      }
    }
    return getSchema()
  })

  const createCredDef = async (schemaId: string): Promise<string> => {
    // wait for schema to be found on the ledger
    // note: in real applications the schema would exist already
    await waitForSchema(schemaId)

    // Create cred def for the provided tag
    console.log(`Creating cred def for schema ID ${schemaId} and tag ${tag}`)
    const msg = new agencyv1.CredDefCreate()
    msg.setSchemaid(schemaId)
    msg.setTag(tag)

    const res = await agentClient.createCredDef(msg)
    console.log(`Cred def created ${res.getId()}`)

    return res.getId()
  }

  const prepareIssuing = async (): Promise<string> => {
    // A dummy schema name
    // Note: creation of schema may fail, if it already exists
    // If this happens, pick a new unique schema name or version and retry
    const schemaName = 'email'
    console.log(`Creating schema ${schemaName}`)

    const schemaMsg = new agencyv1.SchemaCreate()
    schemaMsg.setName(schemaName)
    schemaMsg.setVersion('1.0')
    // List of dummy attributes
    schemaMsg.setAttributesList(['email'])

    try {
      const schemaId = (await agentClient.createSchema(schemaMsg)).getId()
      return await createCredDef(schemaId)
    } catch (err) {
      console.log(`Schema creation failed. Are you trying to recreate an existing schema? ` +
        `Pick an unique schema name or version instead.`)
      process.exit(1)
    }
  }

  // We store the cred def id to a text file
  const credDefIdFilePath = 'CRED_DEF_ID'
  const credDefCreated = existsSync(credDefIdFilePath)
  // Skip cred def creation if it is already created
  const credDefId = credDefCreated ? readFileSync(credDefIdFilePath).toString() : await prepareIssuing()
  if (!credDefCreated) {
    // Store id in order to avoid unnecessary creations
    writeFileSync(credDefIdFilePath, credDefId)
  }
  console.log(`Credential definition available: ${credDefId}`)

  return credDefId
}