<div class="row">
  <div class="col-12">
    <div class="x-panel">
      <div class="x-title">
        <h2>我的Agenda列表</h2>
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
