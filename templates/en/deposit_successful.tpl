<div class="container">
  <div class="header">
    <div class="header-logo">
      <img src="{{ .logo }}">
    </div>
  </div>
  <div class="content">
    <h2 class="title">{{ .record.currency | upcase }} Deposit Successful</h2>
    <p>
      Hi {{ .record.user.email }}
      <br>
      Your deposit of <span class="value">{{ .record.amount }} {{ .record.currency | upcase }}</span> is now available in your SafeTrade account. Login in to check your balance.
      <br>
      <br>
      <br>
      Security Tips:
      <br>
      * Never give your password to anyone.
      <br>
      * Never call any phone number for someone claiming to be Safetrade Support.
      <br>
      * Never send any money to anyone claiming to be a member of Safetrade team.
      <br>
      * Enable Google Two Factor Authentication.
      <br>
      * Bookmark safe.trade and use <a href="#">{{ .record.domain }}/en/official-verification</a> to verify the domain you're visiting.
      <br>
      <br>

      If you don't recognize this activity,please contact our customer support immediately at: <a href="#">{{ .record.domain }}/en/support</a>.
      <br>
      <br>
      <br>

      Safetrade Team
      <br>
      This is an automated message, please do not reply.
    </p>
  </div>
  <p class="footer">
    © 2021 Safetrade All Rights Reserved.
  </p>
</div>
