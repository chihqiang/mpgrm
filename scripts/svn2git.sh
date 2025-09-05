#!/bin/bash
set -euo pipefail

# -----------------------------------------------------------------------------
# 脚本名称: svn2git.sh
# 功能描述: 将标准 SVN 仓库（trunk/branches/tags 结构）迁移为 Git 仓库，
#           自动生成 tag，推送所有分支和 tag（远程不存在的）。
#
# 使用说明:
#   - 支持通过环境变量传入参数：SVN_REPO、GIT_REMOTE_URL、TARGET_DIR（可选）
#   - 支持自定义标签路径：SVN_TAGS_PATH（默认tags）
#   - 支持 dry-run: DRY_RUN=1
#   - 依赖 git-svn 和 subversion（apt/yum/brew 可安装）
#
# 示例:
#   export SVN_REPO="svn://svn.code.sf.net/p/sshpass/code"
#   export GIT_REMOTE_URL="git@github.com:chihqiang/sshpass.git"
#   # 如需自定义标签路径（如releases）：
#   # export SVN_TAGS_PATH="releases"
#   bash scripts/svn2git.sh
# -----------------------------------------------------------------------------

# 配置参数
SVN_REPO="${SVN_REPO:-}"
GIT_REMOTE_URL="${GIT_REMOTE_URL:-}"
TARGET_DIR="${TARGET_DIR:-".svn2git"}"
SVN_TAGS_PATH="${SVN_TAGS_PATH:-"tags"}"  # 允许自定义SVN标签路径
DRY_RUN="${DRY_RUN:-0}"  # 1 表示仅打印操作

# 颜色和格式化定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # 无颜色

# 日志函数
info() {
    echo -e "${BLUE}INFO: $*${NC}"
}

success() {
    echo -e "${GREEN}SUCCESS: $*${NC}"
}

warning() {
    echo -e "${YELLOW}WARNING: $*${NC}"
}

error() {
    echo -e "${RED}ERROR: $*${NC}"
}

# -----------------------------
# 1. 检查依赖
# -----------------------------
check_dependencies() {
    info "Checking required dependencies..."

    # 检查 git 和 svn
    for cmd in git svn; do
        if ! command -v "$cmd" &>/dev/null; then
            error "$cmd not detected, please install it first."
            echo "Ubuntu/Debian: sudo apt install git subversion git-svn"
            echo "CentOS/RHEL : sudo yum install git subversion git-svn"
            echo "macOS       : brew install git svn"
            exit 1
        fi
    done

    # 检查 git-svn 是否可用
    if ! git svn --version &>/dev/null; then
        error "git-svn plugin not available. Please install git-svn."
        exit 1
    fi
    success "All required dependencies are installed"
}

# -----------------------------
# 2. 验证仓库地址
# -----------------------------
validate_repository_urls() {
    info "Validating repository URLs..."

    if [ -z "$SVN_REPO" ]; then
        read -p "Please enter SVN repository URL (SVN_REPO): " SVN_REPO
        [ -z "$SVN_REPO" ] && { error "SVN_REPO not provided, exiting."; exit 1; }
    fi

    if [ -z "$GIT_REMOTE_URL" ]; then
        read -p "Please enter Git remote repository URL (GIT_REMOTE_URL): " GIT_REMOTE_URL
        [ -z "$GIT_REMOTE_URL" ] && { error "GIT_REMOTE_URL not provided, exiting."; exit 1; }
    fi

    # 验证SVN仓库是否可访问
    info "Testing SVN repository connection..."
    if ! svn ls "$SVN_REPO" >/dev/null 2>&1; then
        error "Cannot access SVN repository: $SVN_REPO, please check if the URL is correct"
        exit 1
    fi

    # 显示使用的标签路径
    info "Using SVN tag path: $SVN_TAGS_PATH"
    info "SVN_REPO: $SVN_REPO"
    info "GIT_REMOTE_URL: $GIT_REMOTE_URL"
}

# -----------------------------
# 3. 克隆 SVN 仓库
# -----------------------------
clone_svn_repository() {
    info "Cloning SVN repository to $TARGET_DIR ..."

    if [ "$DRY_RUN" -eq 0 ]; then
        # 清理目标目录
        if [ -d "$TARGET_DIR" ]; then
            warning "Target directory $TARGET_DIR already exists, it will be deleted"
            rm -rf "$TARGET_DIR" || {
                error "Cannot delete directory $TARGET_DIR, please check permissions"
                exit 1
            }
        fi

        # 执行克隆 - 使用自定义标签路径
        info "Cloning with custom layout: trunk=trunk, branches=branches, tags=$SVN_TAGS_PATH"
        if ! git svn clone \
            --trunk=trunk \
            --branches=branches \
            --tags="$SVN_TAGS_PATH" \
            --no-metadata \
            "$SVN_REPO" \
            "$TARGET_DIR"; then
            error "Failed to clone SVN repository"
            exit 1
        fi

        # 进入工作目录
        if ! cd "$TARGET_DIR"; then
            error "Cannot enter directory $TARGET_DIR"
            exit 1
        fi
    else
        info "Simulating clone command: git svn clone --trunk=trunk --branches=branches --tags=$SVN_TAGS_PATH --no-metadata \"$SVN_REPO\" \"$TARGET_DIR\""
        info "Simulating entering directory: cd $TARGET_DIR"
    fi
}

# -----------------------------
# 4. 转换 SVN tags 为 Git tags
# -----------------------------
convert_svn_tags_to_git() {
    info "Converting SVN tags to Git tags ..."

    # 显示所有检测到的SVN远程标签（调试用）
    info "Detecting SVN remote tags, executing: git for-each-ref --format='%(refname:short)' refs/remotes/origin/tags"
    # 尝试直接使用完整路径查找标签
    local svn_tags=$(git for-each-ref --format='%(refname:short)' "refs/remotes/origin/tags" 2>/dev/null || true)

    # 如果没有找到，尝试使用配置的标签路径
    # shellcheck disable=SC2126
    if [ -z "$svn_tags" ] || [ "$(echo "$svn_tags" | grep -v '^$' | wc -l | tr -d ' ')" -eq 0 ]; then
        info "No tags found in refs/remotes/origin/tags, trying configured path: refs/remotes/$SVN_TAGS_PATH"
        svn_tags=$(git for-each-ref --format='%(refname:short)' "refs/remotes/$SVN_TAGS_PATH" 2>/dev/null || true)
    fi

    # 显示检测到的标签数量
    # shellcheck disable=SC2155
    # shellcheck disable=SC2126
    local tag_count=$(echo "$svn_tags" | grep -v '^$' | wc -l | tr -d ' ')
    info "Detected $tag_count SVN tags in total"

    # 显示原始标签列表（用于调试）
    info "Raw tag list:\n$svn_tags"

    # 如果没有检测到标签，显示详细调试信息
    if [ -z "$svn_tags" ] || [ "$tag_count" -eq 0 ]; then
        error "No SVN tags found, this might be caused by path issues"
        error "Currently configured tag path: $SVN_TAGS_PATH"
        error "Try these debug commands to check actual structure:"
        error "svn ls $SVN_REPO/$SVN_TAGS_PATH"
        error "In migration directory: git branch -r | grep 'tags'"
        # 列出所有远程分支供调试
        info "All remote branches:\n$(git branch -r 2>/dev/null || true)"
        return
    fi

    # 显示检测到的标签列表
    info "Detected SVN tags list:"
    echo "$svn_tags" | grep -v '^$' | while read -r tag; do
        echo "  - $tag"
    done

    # 转换标签
    echo "$svn_tags" | grep -v '^$' | while read -r tag; do
        # 提取纯标签名（去掉前缀）
        # 尝试多种可能的前缀格式
        local tagname
        if [[ "$tag" == origin/tags/* ]]; then
            tagname="${tag#origin/tags/}"
        elif [[ "$tag" == $SVN_TAGS_PATH/* ]]; then
            tagname="${tag#$SVN_TAGS_PATH/}"
        else
            # 如果没有匹配的前缀，使用标签全名
            tagname="$tag"
        fi

        # 确保标签名不为空
        if [ -z "$tagname" ]; then
            warning "Invalid tag name (empty), skipping: $tag"
            continue
        fi

        # 检查标签是否已存在
        if git rev-parse "refs/tags/$tagname" >/dev/null 2>&1; then
            info "Tag already exists, skipping: $tagname"
            continue
        fi

        info "Creating tag: $tagname (from $tag)"
        if [ "$DRY_RUN" -eq 0 ]; then
            # 显示完整的标签创建命令（调试用）
            info "Executing: git tag -a \"$tagname\" -m \"Tag imported from SVN: $tagname\" \"refs/remotes/$tag\""
            if git tag -a "$tagname" -m "Tag imported from SVN: $tagname" "refs/remotes/$tag"; then
                success "Successfully created tag: $tagname"
            else
                error "Failed to create tag $tagname, please check error message"
            fi
        fi
    done

    # 显示转换后的本地标签数量
    if [ "$DRY_RUN" -eq 0 ]; then
        # shellcheck disable=SC2155
        local git_tag_count=$(git tag | wc -l | tr -d ' ')
        info "Conversion completed, created $git_tag_count local Git tags"
        # 显示创建的标签列表
        info "Created Git tags list:\n$(git tag 2>/dev/null || true)"
        if [ "$git_tag_count" -eq 0 ]; then
            error "No Git tags were successfully created, this might be due to tag path parsing issues or insufficient permissions"
            error "Suggest manually checking SVN tag structure: svn ls $SVN_REPO/$SVN_TAGS_PATH"
            error "And try creating tags manually: git tag -a <tag-name> -m <message> <commit-hash>"
        fi
    fi
}

# -----------------------------
# 5. 清理远程分支和 SVN 元数据
# -----------------------------
cleanup_svn_metadata() {
    info "Cleaning up remote branches and SVN metadata ..."

    if [ "$DRY_RUN" -eq 0 ]; then
        # 清理SVN元数据
        rm -rf .git/svn .git/logs/refs/remotes 2>/dev/null || true

        # 删除远程跟踪分支
        # shellcheck disable=SC2155
        local remote_branches=$(git branch -r 2>/dev/null || true)
        if [ -n "$remote_branches" ]; then
            for b in $remote_branches; do
                info "Deleting remote tracking branch: $b"
                git branch -rd "$b" 2>/dev/null || true
            done
        else
            info "No remote tracking branches to clean up"
        fi
    fi
}

# -----------------------------
# 6. 添加远程仓库
# -----------------------------
add_remote_repository() {
    info "Configuring Git remote repository: $GIT_REMOTE_URL"

    if [ "$DRY_RUN" -eq 0 ]; then
        # 检查是否已存在origin远程
        if git remote | grep -q "^origin$"; then
            info "Updating existing origin remote URL"
            git remote set-url origin "$GIT_REMOTE_URL"
        else
            info "Adding origin remote repository"
            git remote add origin "$GIT_REMOTE_URL"
        fi

        # 验证远程连接
        info "Validating remote repository connection..."
        if ! git ls-remote --exit-code origin >/dev/null 2>&1; then
            error "Cannot connect to remote repository $GIT_REMOTE_URL, please check URL and permissions"
            exit 1
        fi
    fi
}

# -----------------------------
# 7. 推送所有本地分支
# -----------------------------
push_local_branches() {
    info "Pushing all local branches..."

    # shellcheck disable=SC2155
    local local_branches=$(git for-each-ref --format='%(refname:short)' refs/heads/ 2>/dev/null || true)

    if [ -z "$local_branches" ]; then
        warning "No local branches found"
        return
    fi

    for branch in $local_branches; do
        info "Pushing branch: $branch"
        if [ "$DRY_RUN" -eq 0 ]; then
            if git push origin "$branch"; then
                success "Successfully pushed branch: $branch"
            else
                error "Failed to push branch $branch"
            fi
        fi
    done
}

# -----------------------------
# 8. 推送缺失的标签（只推送远程不存在的标签）
# -----------------------------
push_missing_tags() {
    info "Pushing missing tags (only push tags not existing remotely)..."

    # 获取本地标签列表
    info "Getting local tag list: git tag"
    # 先显示git tag的原始输出
    info "git tag raw output:\n$(git tag 2>&1)"
    # 然后获取并处理标签列表
    # shellcheck disable=SC2155
    local local_tags=$(git tag | sort | uniq 2>/dev/null || true)
    info "Processed local tag list:\n$local_tags"

    # 获取远程标签列表
    info "Getting remote tag list: git ls-remote --tags origin"
    # 先显示原始输出
    info "git ls-remote raw output:\n$(git ls-remote --tags origin 2>&1)"
    # 然后获取并处理远程标签列表
    # shellcheck disable=SC2155
    local remote_tags=$(git ls-remote --tags origin | awk '{print $2}' | sed 's|refs/tags/||' | sort | uniq 2>/dev/null || true)
    info "Processed remote tag list:\n$remote_tags"

    if [ -z "$local_tags" ]; then
        error "No local tags found, nothing to push"
        error "Please check if tag conversion step was successful"
        error "Possible reasons: tag creation failed, insufficient permissions, or Git configuration issues"
        # 显示更多调试信息
        info "Current Git repository status:\n$(git status 2>&1)"
        info "Git tag list command exit code: $?"
        return
    fi

    # 显示本地和远程标签数量
    # shellcheck disable=SC2155
    local local_tag_count=$(echo "$local_tags" | grep -v '^$' | wc -l | tr -d ' ')
    # shellcheck disable=SC2155
    local remote_tag_count=$(echo "$remote_tags" | grep -v '^$' | wc -l | tr -d ' ')
    info "Local tag count: $local_tag_count"
    info "Remote tag count: $remote_tag_count"

    # 找出远程不存在的本地标签
    # shellcheck disable=SC2155
    local missing_tags=$(comm -23 <(echo "$local_tags") <(echo "$remote_tags"))
    
    if [ -z "$missing_tags" ]; then
        info "All local tags already exist in remote repository, no need to push"
        return
    fi

    # 显示将要推送的标签列表
    info "Tags to be pushed (not existing remotely):"
    echo "$missing_tags" | while read -r tag; do
        echo "  - $tag"
    done

    # 推送标签
    local tags_pushed=0
    local tags_failed=0
    
    echo "$missing_tags" | while read -r tag; do
        info "Pushing tag: $tag"
        if [ "$DRY_RUN" -eq 0 ]; then
            # 显示完整的推送命令（调试用）
            info "Executing: git push origin \"refs/tags/$tag\""
            if git push origin "refs/tags/$tag"; then
                success "Successfully pushed tag: $tag"
                ((tags_pushed++))
            else
                error "Failed to push tag $tag, trying force push..."
                info "Executing: git push -f origin \"refs/tags/$tag\""
                if git push -f origin "refs/tags/$tag"; then
                    success "Successfully force pushed tag: $tag"
                    ((tags_pushed++))
                else
                    error "Failed to force push tag $tag"
                    error "Error message: $(git push -f origin "refs/tags/$tag" 2>&1)"
                    ((tags_failed++))
                fi
            fi
        else
            # Dry-run模式下计数
            ((tags_pushed++))
        fi
    done

    # 显示推送统计
    info "Tag push statistics: $tags_pushed succeeded, $tags_failed failed, $((local_tag_count - tags_pushed - tags_failed)) skipped"
    # shellcheck disable=SC2031
    if [ "$tags_failed" -gt 0 ]; then
        warning "Some tags failed to push,suggestion check network connection and permissions"
        warning "Try pushing failed tags manually: git push -f origin --tags"
    elif [ "$tags_pushed" -eq 0 ]; then
        info "No new tags were pushed, all local tags already exist in remote"
    fi
}

# -----------------------------
# 主程序
# -----------------------------
main() {
    echo "============================================="
    echo "          SVN to Git repository migration tool             "
    echo "============================================="
    [ "$DRY_RUN" -eq 1 ] && warning "Runs in DRY-RUN mode and does not perform any actual operations"
    echo

    check_dependencies
    echo

    validate_repository_urls
    echo

    clone_svn_repository
    echo

    convert_svn_tags_to_git
    echo

    cleanup_svn_metadata
    echo

    add_remote_repository
    echo

    push_local_branches
    echo

    push_missing_tags
    echo

    success "Migration process completed！"
    [ "$DRY_RUN" -eq 1 ] && warning "Note: This is DRY-RUN mode, no actual operation is performed"
}

# 启动主程序
main
