-- =========================
-- USERS TABLE
-- =========================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- ROLES TABLE
-- =========================
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- =========================
-- PERMISSIONS TABLE
-- =========================
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- =========================
-- USER_ROLES (M:N)
-- =========================
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,

    PRIMARY KEY (user_id, role_id),

    CONSTRAINT fk_user_roles_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_roles_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE
);

-- =========================
-- ROLE_PERMISSIONS (M:N)
-- =========================
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,

    PRIMARY KEY (role_id, permission_id),

    CONSTRAINT fk_role_permissions_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_role_permissions_permission
        FOREIGN KEY (permission_id)
        REFERENCES permissions(id)
        ON DELETE CASCADE
);

-- =========================
-- INDEXES (performance for IAM lookups)
-- =========================
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id
    ON user_roles(user_id);

CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id
    ON role_permissions(role_id);