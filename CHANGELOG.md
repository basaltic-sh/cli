# Changelog

## 0.30.0

- Choose the exact instance-pool members to remove with `instance-pool update --scale-in-instances`; use `--idempotency-key` to retry a scale-down safely.
- Set or clear a member's scale-down protection with `instance-pool update-instance --scale-in-protected=true` or `--scale-in-protected=false`.
- Quota list responses now include quota policies alongside quota usage.
