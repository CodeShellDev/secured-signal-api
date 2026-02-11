package common

type contextKey string

const TokenKey contextKey = "token"
const IsAuthKey contextKey = "isAuthenticated"

const LoggerKey contextKey = "logger"

const TrustedClientKey contextKey = "isClientTrusted"

const TrustedProxyKey contextKey = "isProxyTrusted"
const ClientIPKey contextKey = "clientIP"
const OriginURLKey contextKey = "originURL"