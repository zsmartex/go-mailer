<!doctype html>
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <style>
      *, *::before, *::after {
        box-sizing: border-box;
      }
      strong {
        margin: 0 4px;
      }
      .outer-container {
        background-color: #F8F8F8;
        font-family: "Roboto";
        font-size: 14px;
        line-height: 1.4;
        margin: 0;
        padding: 26px 0;
      }
      .container {
        max-width: 660px;
        margin: 0 auto;
      }
      .header {
        width: 100%;
        border-top-left-radius: 8px;
        border-top-right-radius: 8px;
        background-color: #213150;
      }
      .header-logo {
        text-align: center;
        padding: 12px 0;
      }
      .header-logo img {
        height: 75px;
      }
      .content {
        background-color: #fff;
        padding: 2em 1em;
      }
      .content .title {
        line-height: 1.5;
        font-size: 1.5em;
        font-weight: bold;
        margin: 0;
      }
      .content p {
        font-size: 14px;
        font-weight: normal;
        margin: 16px 0 0 0;
        color: #474d57;
      }
      .content a {
        color: #0095ff;
        text-decoration: none;
      }
      .confirm_code {
        font-size: 18px;
        font-weight: bold;
        color: #0095ff;
      }
      .value {
        font-weight: 500;
      }
      table .table_title {
        font-weight: normal;
        color: #76808f;
        font-size: 14px;
        padding-right: 8px;
        padding-left: 8px;
      }
      .footer {
        text-align: center;
        line-height: 1.5;
        padding: 15px;
        font-size: 1em;
        color: #999999;
      }
    </style>
  </head>

  <body>
    <div class="outer-container">
      {{ .Body }}
    </div>
  </body>
</html>
