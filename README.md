# mpgrm - Multi-Platform Git Repository Manager

A lightweight tool for multi-platform Git repository migration, supporting selective transfer of branches and tags between different code hosting platforms.

## 🔥 Key Features

- **Cross-Platform Repository Migration**: Seamlessly migrate repositories between different Git code hosting platforms
- **Selective Migration**: Support for migrating only specific branches and tags
- **Release Management**: Upload, download, create, and sync releases and attachments
- **Bulk Operations**: Support for batch repository synchronization at the organization or user level
- **Multi-Platform Support**: Compatible with GitHub, Gitee, Gitea, and cnb.cool

## 🚀 Installation Guide

### Install from Source

- Ensure you have Go 1.23+ installed
- Clone the repository and compile

```bash
git clone https://github.com/chihqiang/mpgrm.git
cd mpgrm && make build
```

The compiled executable will be in the project root directory

### Environment Configuration

Copy the environment variable example file

```bash
cp .env.example .env
```

Edit the .env file and add the corresponding authentication information based on your platform

```env
# CNB platform username (CNB does not have personal repositories; only organizations exist. 
# This value is used as the "personal repository" identifier)
CNB_USERNAME="cnb"
# CNB platform access token
CNB_TOKEN="token_cnb_abcdef123456"

# Gitea platform username
GITEA_USERNAME="gitea"
# Gitea platform access token
GITEA_TOKEN="a1b2c3d4e5f6g7h8i9j0klmnopqrstuvwx"

# Gitee platform username
GITEE_USERNAME="gitee"
# Gitee platform access token
GITEE_TOKEN="z9y8x7w6v5u4t3s2r1q0ponmlkjihgfedcba"

# GitHub username
GITHUB_USERNAME="github"
# GitHub personal access token
GITHUB_TOKEN="ghp_FAKE1234567890abcdefABCDEFabcdef"

```

## 💻 Usage Examples

### credential

Manage repository credentials with ease and push anytime

~~~
# Basic usage
# Specify the repository URL directly via the command line. 
# The default behavior (whether to read from environment variables) is controlled by --use-env.
mpgrm credential show --repo https://github.com/username/source-repo.git

# Specify the repository URL and explicitly disable reading from environment variables.
# --use-env=false means the command will only use the URL provided via --repo.
mpgrm credential show --use-env=false --repo https://github.com/username/source-repo.git
~~~

### Push Repository (push)

Push specified branches and tags from source repository to target repository

```bash
# Basic usage
mpgrm push --repo https://github.com/username/source-repo.git --target-repo https://gitee.com/username/target-repo.git

# Selective push of specific branches and tags
mpgrm push --repo https://github.com/username/source-repo.git --target-repo https://gitee.com/username/target-repo.git --branches main,develop --tags v1.0.0,v1.1.0

# Specify workspace directory
mpgrm push --repo https://github.com/username/source-repo.git --target-repo https://gitee.com/username/target-repo.git --workspace /path/to/workspace
```

### Manage Releases (releases)

#### Upload Release Files

```bash
mpgrm releases upload --repo https://github.com/username/repo.git --tags v1.0.0 --files path/to/file1.zip,path/to/file2.tar.gz
```

#### Download Release Files

```bash
mpgrm releases download --repo https://github.com/username/repo.git --tags v1.0.0,v1.1.0
```

#### Create Releases for All Tags

```bash
mpgrm releases create --repo https://github.com/username/repo.git
```

#### Sync Releases

```bash
mpgrm releases sync --repo https://github.com/username/source-repo.git --target-repo https://gitee.com/username/target-repo.git --tags v1.0.0,v1.1.0
```

### Manage Repositories (repo)

#### List Repositories

```bash
# List all repositories of an organization
mpgrm repo list --repo https://github.com/organization/

# List all repositories of a user
mpgrm repo list --repo https://github.com
```

####  Clone Repositories

~~~bash
# List and clone all repositories of an organization
mpgrm repo clone --repo https://github.com/<organization>/

# List and clone all repositories of a user
mpgrm repo clone --repo https://github.com/<username>
~~~

#### Sync Repositories

```bash
# Sync all repositories of an organization
mpgrm repo sync --repo https://github.com/organization/ --target-repo https://gitee.com/organization/

# Sync all repositories of a user
mpgrm repo sync --repo https://github.com --target-repo https://gitee.com
```

### repo & target-repo Usage Guide

CloneURL determines whether it points to an **organization** or a **repository** based on the trailing character.

#### Organization URL (ends with `/`)

- Indicates a organization path.
- The repository name should be appended to generate a full repository URL.
- Currently, URLs without a trailing `/` (e.g., `https://github.com/org`) are NOT supported.
- Examples:

```yaml
https://github.com/org/ -> org
https://cnb.cool/org1/ -> org1
https://cnb.cool/org1/child1/ -> org1/child1
https://cnb.cool/org1/child1/child2/ -> org1/child1/child2
```

## Personal Homepage or Root URL

- URLs that point to the root of a domain or user homepage.
- Usually just for reference; no repository is implied.

- Examples:

```yml
https://github.com/
https://github.com
https://cnb.cool
https://cnb.cool/
```

#### Repository URL (does not end with `/` or ends with `.git`)

```yml
https://github.com/org/repo.git
https://cnb.cool/org1/child1
```

> ⚠️ Note:  
> The **trailing `/`** is the key to distinguish between organization URLs and repository URLs.  
> Organization URLs require appending a repository name, whereas repository URLs are already complete.

## 🤝 Contribution Guide

1. Fork this repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request


## 📧 Contact

If you have any questions or suggestions, please feel free to submit issues or contact the project maintainers.
