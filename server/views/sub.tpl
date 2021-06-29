<div class="row">
  <div class="col-12">
    <div class="x-panel">
      <div class="x-title">
        <h2>我的Agenda列表</h2>
        <div class="x-cmd">
          <a role="button" class="btn btn-primary" href="/icsgen">生成日历</a>
        </div>
      </div>
      <div class="x-content">
        <table>
          <tbody>
            {{range .files}}
            <tr>
              <td>{{.Name}}</td>
              <td>{{.Status}}</td>
            </tr>
            {{end}}
          </tbody>
        </table>
      </div>
  </div>
</div>
