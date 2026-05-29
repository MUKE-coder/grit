# Grit Motors — Default Login

A fresh deployment ships with exactly **one** account so you can log in
and create everything else through the UI.

| Role  | Email                       | Password      |
|-------|-----------------------------|---------------|
| Admin | admin@grit.demo      | password123   |

The seeder also creates:
- Business: **Grit Motors**
- Branch: **Main Branch** (default)

## First steps after login

1. Settings → change the admin password.
2. Branches → add additional branches (Busega, Bweyogerere, Gayaza, etc.).
3. Staff → invite managers, loan officers, cashiers.
4. Loan Products → define your loan templates.
5. Motorcycles → start adding inventory.

## Roles

- **admin** — full access across both workspaces
- **manager** — operations + reports for their branch
- **loan_officer** — borrowers, loans, repayments (no spares POS)
- **accountant** — read-only finance reports across segments
- **cashier** — POS + cash sales + collect repayments (pending verification)
- **stock_clerk** — receive stock only
