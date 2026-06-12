import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        TOTP — Time-Based One-Time Password — is the 6-digit code your
        authenticator app generates (1Password, Authy, Google Authenticator).
        Grit ships TOTP 2FA with encrypted seed storage. Users opt in;
        admins can enforce it.
      </p>

      <h2>The setup flow</h2>
      <CodeBlock
        language="text"
        code={`User clicks "Enable 2FA" in settings
  → POST /api/auth/totp/setup
       returns:
         { secret: "JBSWY3DPEHPK3PXP", qr_code: "data:image/png;base64..." }
  → User scans QR with their app
  → User types the 6-digit code from the app
  → POST /api/auth/totp/enable { code: "123456" }
       Grit verifies the code, marks 2FA enabled on the user`}
      />

      <h2>The login flow with TOTP enabled</h2>
      <CodeBlock
        language="text"
        code={`User submits email + password to /api/auth/login
  → If 2FA enabled: response is { totp_challenge: <one-time-token> }
                    NOT a full access token yet
User opens their authenticator app, types the 6-digit code
  → POST /api/auth/totp/verify { challenge: <token>, code: "123456" }
       → Grit verifies, returns the real access + refresh tokens`}
      />

      <h2>Encrypted seed storage</h2>
      <p>
        TOTP&apos;s security rests on the seed (the shared secret stored on the
        server). If a SQL injection or DB dump leaks it, the attacker can
        generate valid codes forever. Grit&apos;s defence: <strong>seeds are
        encrypted at rest</strong> with AES-GCM, keyed off{' '}
        <code>JWT_SECRET</code>.
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/totp/store.go (excerpt)"
        code={`func (s *Store) StoreSeed(userID, seed string) error {
    nonce := make([]byte, s.aead.NonceSize())
    rand.Read(nonce)
    encrypted := s.aead.Seal(nonce, nonce, []byte(seed), nil)
    return s.db.Save(&UserTOTP{
        UserID:        userID,
        EncryptedSeed: encrypted,
    }).Error
}`}
      />
      <p>
        A SQL injection that leaks the table gets ciphertext — useless
        without the env-var key.
      </p>

      <TipBox tone="success">
        Recovery codes — Grit also generates 8 one-time recovery codes when
        TOTP is enabled. Each can be used once to log in if the
        authenticator is lost. Show them once, never again, force the user
        to save them.
      </TipBox>

      <h2>The trusted-device flow</h2>
      <p>
        Re-prompting for a TOTP code on every login is painful. Grit&apos;s
        flow: on successful TOTP verify, the user can tick &quot;trust this
        device for 30 days&quot;. Grit issues a device-fingerprint cookie. Next
        login on the same device — fingerprint matches, TOTP is skipped.
      </p>

      <h2>What can break</h2>
      <ul>
        <li>
          <strong>Clock drift</strong> — TOTP codes are time-windowed. If the
          server&apos;s clock is way off, codes will never verify. Grit allows ±1
          step (30 seconds) of drift by default.
        </li>
        <li>
          <strong>Wrong seed encoding</strong> — base32 vs base64 mismatches
          are a classic bug. Grit&apos;s setup endpoint returns base32 (what
          authenticator apps expect).
        </li>
        <li>
          <strong>Encrypted seed migration</strong> — if you rotate{' '}
          <code>JWT_SECRET</code>, existing TOTP seeds can&apos;t be decrypted.
          Plan for re-enrolment if you rotate.
        </li>
      </ul>

      <KnowledgeCheck
        question="A SQL injection leaks the user_totps table. The attacker has every encrypted seed. What can they do?"
        choices={[
          {
            label: 'Generate valid TOTP codes for every user immediately.',
            feedback:
              "Wrong — the seeds are encrypted. Without the decryption key (derived from JWT_SECRET, which lives in env, not the DB), the leaked ciphertext is useless.",
          },
          {
            label: 'Nothing — the seeds are AES-GCM encrypted with a key the attacker doesn\'t have.',
            correct: true,
            feedback:
              "Right — that's why encrypted seed storage is a meaningful defence. The Defender's Handbook lesson on TOTP seed theft covered this exact attack and exactly this defence.",
          },
          {
            label: 'Brute-force the encryption — AES-GCM is breakable.',
            feedback:
              "Wrong — AES-GCM with a properly random 256-bit key is not currently brute-forceable. They'd need the JWT_SECRET (or compromise the running server) to decrypt.",
          },
          {
            label: 'They can bypass TOTP entirely by sending an empty seed.',
            feedback:
              "Wrong — the verifier always compares against the stored (encrypted) seed; there's no 'empty seed' path.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Enable TOTP on your own account:</p>
            <ol>
              <li>
                Log in via the API. Grab the access token.
              </li>
              <li>
                <code>POST /api/auth/totp/setup</code> with the bearer token.
                Save the QR data URL.
              </li>
              <li>
                Open the QR URL in your browser (paste it into the address
                bar) — your browser renders the base64 image.
              </li>
              <li>Scan with 1Password / Authy / Google Authenticator.</li>
              <li>
                Type the 6-digit code to{' '}
                <code>POST /api/auth/totp/enable</code> with the bearer token.
              </li>
              <li>
                Log out, log back in — you should now see the{' '}
                <code>totp_challenge</code> response instead of tokens.
                Submit the new code to{' '}
                <code>POST /api/auth/totp/verify</code>.
              </li>
            </ol>
            <p>
              Paste your <code>totp_challenge</code> response in{' '}
              <code>notes.md</code> as proof.
            </p>
          </>
        }
        hint={
          <>
            If the QR doesn&apos;t render, the response value already starts with{' '}
            <code>data:image/png;base64,</code>; paste the whole string into
            the address bar.
          </>
        }
        solution={
          <>
            <p>
              The totp_challenge response shape:
            </p>
            <CodeBlock
              language="json"
              code={`{
  "totp_challenge": "<one-time short-lived token>",
  "message": "Enter the 6-digit code from your authenticator"
}`}
            />
            <p>
              You&apos;ve now enrolled in 2FA, simulated a login flow that
              requires it, and have working recovery codes saved.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson of the chapter — <strong>RBAC + invitations</strong>.
        Roles, ownership, team invites. The pattern every multi-tenant app
        needs.
      </p>
    </>
  )
}
