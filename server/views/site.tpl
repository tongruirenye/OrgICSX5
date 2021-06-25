<!doctype html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css" rel="stylesheet" crossorigin="anonymous">
    <style type="text/css">
     body {
       background: #F7F7F7;
     }
     .x-panel {
       position: relative;
       width: 100%;
       margin-bottom: 10px;
       padding: 10px 17px;
       background: #fff;
       border: 1px solid #E6E9ED;
       opacity: 1;
       transition: all .2s ease;
     }
     .x-title {
       border-bottom: 2px solid #E6E9ED;
       padding: 1px 5px 6px;
       margin-bottom: 10px;
     }
    </style>
    <title>我的TODO</title>
  </head>
  <body>
    <div class="container">
      <nav class="navbar navbar-expand-sm navbar-dark bg-dark">
        <div class="container">
          <a class="navbar-brand" href="#">
            <img src="/docs/5.0/assets/brand/bootstrap-logo.svg" alt="" width="30" height="24">
          </a>
          <div class="collapse navbar-collapse">
            <ul class="navbar-nav me-auto">
              <li class="navbar-item">
                <a class="nav-link active" href="#">Today</a>
              </li>
              <li class="navbar-item">
                <a class="nav-link" href="/sub">订阅</a>
              </li>
            </ul>
            <ul class="navbar-nav">
              <li class="nav-item dropdown">
                <a class="nav-link dropdown-toggle" id="loginDropDown" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">{{.username}}</a>
                <ul class="dropdown-menu" aria-labelledby="loginDropDown">
                  <li>
                    <a class="dropdown-item" href="#">登出</a>
                  </li>
                </ul>
              </li>
            </ul>
          </div>
        </div>
      </nav>
      {{.LayoutContent}}
    </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/js/bootstrap.bundle.min.js" crossorigin="anonymous"></script>
  </body>
</html>
