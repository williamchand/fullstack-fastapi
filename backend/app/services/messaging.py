# app/services/messaging.py
def send_email(to_email: str, subject: str, body: str) -> None:
    # TODO: integrate SendGrid / Mailgun / SMTP
    print(f"[EMAIL] to={to_email} subject={subject} body={body}")

def send_sms(to_phone: str, body: str) -> None:
    # TODO: integrate Twilio / Vonage / etc.
    print(f"[SMS] to={to_phone} body={body}")
