# 1Password Setup for Mondrian Notion Integration

## 🔐 Step-by-Step Setup

### 1. Install 1Password CLI
```bash
brew install 1password-cli
```

### 2. Sign in to 1Password
```bash
op signin
# Follow the prompts to authenticate
```

### 3. Create Secrets in 1Password

**Option A: Using 1Password App (Recommended)**
1. Open 1Password app
2. Create a new **API Credential** item called "Notion API Token"
   - Title: `Notion API Token`
   - Username: `notion-integration`
   - Token: `[your_notion_integration_token]`
   - Website: `https://notion.so`

3. Create another **Database** item called "Mondrian Notion Database"
   - Title: `Mondrian Notion Database`  
   - Database ID: `[your_notion_database_id]`
   - Website: `https://notion.so/mondrian-workspace`

**Option B: Using CLI**
```bash
# Create Notion API token
op item create --category="API Credential" \
  --title="Notion API Token" \
  token="your_integration_token_here"

# Create database ID
op item create --category="Database" \
  --title="Mondrian Notion Database" \
  database_id="your_database_id_here"
```

### 4. Test Secret Retrieval
```bash
# Test reading the secrets
op read "op://Private/Notion API Token/token"
op read "op://Private/Mondrian Notion Database/database_id"
```

### 5. Update Script References
If your items are in a different vault or have different names, update the references in `notion-mcp-bridge.ts`:

```typescript
// Current references (update if needed):
const token = await get1PasswordSecret('op://Private/Notion API Token/token');
const databaseId = await get1PasswordSecret('op://Private/Mondrian Notion Database/database_id');
```

## 🏃‍♂️ Alternative: Run with 1Password Injection

**Option 1: Direct command injection**
```bash
op run --env-file=".env.1p" -- npm run test-notion
```

Create `.env.1p` file:
```
NOTION_TOKEN="op://Private/Notion API Token/token"
NOTION_DATABASE_ID="op://Private/Mondrian Notion Database/database_id"
```

**Option 2: Inline injection**
```bash
op run -- sh -c 'NOTION_TOKEN="op://Private/Notion API Token/token" NOTION_DATABASE_ID="op://Private/Mondrian Notion Database/database_id" npm run test-notion'
```

## 🔍 Getting Your Notion Credentials

### 1. Create Notion Integration
1. Go to https://www.notion.so/my-integrations
2. Click "New integration"
3. Name: `Mondrian MCP Bridge`
4. Associated workspace: Your workspace
5. Copy the **Integration Token** → Save to 1Password

### 2. Get Database ID
1. Open your Mondrian todos database in Notion
2. Copy the URL: `https://notion.so/workspace/DATABASE_ID?v=...`
3. Extract the `DATABASE_ID` part (32 characters)
4. Save to 1Password

### 3. Share Database with Integration
1. In your Notion database, click "Share" 
2. Click "Invite"
3. Search for your integration name: "Mondrian MCP Bridge"
4. Give it "Edit" permissions

## 🧪 Testing

```bash
cd /Users/chrismcconnell/GitHub/mondrian
npm run test-notion
```

Expected output:
```
✅ Successfully retrieved secrets from 1Password
✅ Notion MCP Bridge initialized with data_source_id: datasource_xxx
📋 Found X existing todos
```

## 🚨 Security Notes

- Never commit actual tokens to git
- 1Password CLI requires authentication every ~30 minutes
- Secret references (`op://...`) are safe to commit
- Use `op signin --force` if authentication expires

## 🔧 Troubleshooting

**"command not found: op"**
```bash
brew install 1password-cli
```

**"authentication required"**
```bash
op signin --force
```

**"item not found"**
- Check item names match exactly
- Verify vault name (default: "Private")
- Use `op item list` to see all items

**"insufficient permissions"**
- Make sure integration has access to your Notion database
- Re-share database with integration if needed