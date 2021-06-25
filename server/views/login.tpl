<!doctype html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css" rel="stylesheet" crossorigin="anonymous">

    <style type="text/css">
     .wrapper {
       max-width: 360px;
       margin: 0 auto;
       text-align: center;
     }

     .wrapper input[type="text"],
     .wrapper input[type="password"] {
       margin: 0 0 20px;
     }

     .wrapper h1 {
       margin: 10px 0 30px;
     }
    </style>
    <title>我的TODO</title>
  </head>
  <body>
    <div class="container">
    <div class="wrapper mt-3">
      <form action="/login" method="POST">
        <h1>登陆</h1>
        {{.xsrfdata}}
        {{if .flash.error}}
        <div class="alert alert-warning alert-dismissible fade show" role="alert">
          {{.flash.error}}
          <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
        </div>
        {{end}}
        <input type="text" class="form-control" placeholder="邮箱" required="required" name="username" />
        <input type="password" class="form-control" placeholder="密码" required="required" name="password" />
        <button type="submit" class="btn btn-primary">登陆</button>
      </form>
    </div>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/js/bootstrap.bundle.min.js" crossorigin="anonymous"></script>
  </body>
</html>
