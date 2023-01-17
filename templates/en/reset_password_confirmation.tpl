<div class="container">
  <div class="header">
    <div class="header-logo">
      <img src="{{ .logo }}">
    </div>
  </div>
  <div class="content">
    <h2 class="title">Reset Password Verification Code</h2>
    <p>
      Welcome to ZSmartex
      <br>
      Your reset password verification code is <span class="confirm_code">{{ .record.code }}</span> It was generated at {{ .record.created_at }} and will be valid for {{ .record.expired_in }} minutes.
      <br>
      <br>
      <br>
      Security Tips:
      <br>
      * Never give your password to anyone.
      <br>
      * Never call any phone number for someone claiming to be ZSmartex Support.
      <br>
      * Never send any money to anyone claiming to be a member of ZSmartex team.
      <br>
      * Enable Google Two Factor Authentication.
      <br>
      * Bookmark www.zsmartex.tech and use <a href="#">{{ .record.domain }}/en/official-verification</a> to verify the domain you're visiting.
      <br>
      <br>

      If you don't recognize this activity,please contact our customer support immediately at: <a href="#">{{ .record.domain }}/en/support</a>.
      <br>
      <br>
      <br>

      ZSmartex Team
      <br>
      This is an automated message, please do not reply.
    </p>
  </div>
  <p class="footer">
    © 2021 ZSmartex.tech All Rights Reserved.
  </p>
</div>
