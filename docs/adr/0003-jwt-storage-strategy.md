# 0003 JWT Storage Strategy

## Status

Accepted

## Context

We need to decide where clients store the JWT token. Options: localStorage,
httpOnly cookie, in-memory.

## Decision

Store JWT in localStorage on the frontend. The backend validates the token on
every request via the `Authorization: Bearer` header.

## Consequences

Simple to implement. Token survives page refresh. Vulnerable to XSS,
acceptable for this project scope. For higher security requirements, switch
to httpOnly cookies to prevent JS access.
