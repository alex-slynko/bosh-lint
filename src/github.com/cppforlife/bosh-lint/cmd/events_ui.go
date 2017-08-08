package cmd

/* todo
- pagination
- style?
- api endpoints for:
  - t 8 --json

*/

const eventsUI = `
<!DOCTYPE html>
<html>
<head>
  <script
  src="http://code.jquery.com/jquery-3.2.1.min.js"
  integrity="sha256-hwg4gsxgFZhOsEEamdOYGBf13FyQuiTwlAQgxVSNgt4="
  crossorigin="anonymous"></script>
</head>
<body>

<div id="root"></div>

<div id="canvas-tmpl" class="tmpl">
  <form autocomplete="off">
    <input type="text" name="after" placeholder="after" />
    <input type="text" name="before" placeholder="before" />
    <input type="text" name="event-user" placeholder="user" />
    <input type="text" name="action" placeholder="action" />
    <input type="text" name="object-type" placeholder="obj. type" />
    <input type="text" name="object-name" placeholder="obj. name" />
    <input type="text" name="task" placeholder="task" />
    <input type="text" name="deployment" placeholder="deployment" />
    <input type="text" name="instance" placeholder="instance" />
    <button type="submit">Go</button>
  </form>
</div>

<table id="event-tmpl" class="tmpl">
  <tr data-event-id="{id}">
    <td class="id">{id}</td>
    <td class="time">
      <a href="#" data-query="after" data-value="{time}">{time}</a>
    </td>
    <td class="user">
      <a href="#" data-query="event-user" data-value="{user}">{user}</a>
    </td>
    <td class="action">
      <a href="#" data-query="action" data-value="{action}">{action}</a>
    </td>
    <td class="object_type">
      <a href="#" data-query="object-type" data-value="{object_type}">{object_type}</a>
    </td>
    <td class="object_name">
      <a href="#" data-query="object-name" data-value="{object_name}">{object_name}</a>
    </td>
    <td class="task_id">
      <a href="#" data-query="task" data-value="{task_id}">{task_id}</a>
    </td>
    <td class="deployment">
      <a href="#" data-query="deployment" data-value="{deployment}">{deployment}</a>
    </td>
    <td class="instance">
      <a href="#" data-query="instance" data-value="{instance}">{instance}</a>
    </td>
    <td class="context"><span>{context}</span></td>
    <td class="error"><span>{error}</span></td>
  </tr>
</table>

<table id="no-events-tmpl" class="tmpl">
  <tr><td colspan="11">No matching events</td></tr>
</table>

<table id="error-events-tmpl" class="tmpl">
  <tr><td colspan="11">Error fetching events</td></tr>
</table>

<table id="task-tmpl" class="tmpl">
  <tr>
    <td class="id">{id}</td>
    <td class="started_at">{started_at}</td>
    <td class="last_activity_at">{last_activity_at}</td>
    <td class="user">{user}</td>
    <td class="deployment">{deployment}</td>
    <td class="description">{description}</td>
    <td class="result">{result}</td>
  </tr>
</table>

<table id="no-tasks-tmpl" class="tmpl">
  <tr><td colspan="11">No matching tasks</td></tr>
</table>

<table id="error-tasks-tmpl" class="tmpl">
  <tr><td colspan="11">Error fetching tasks</td></tr>
</table>

<script type="text/javascript">

function CanvasCollection($el) {
  var $canvases = null;

  function setUp() {
    $canvases = newNamedDiv($el, "canvases")
    NewCanvasAddButton(newNamedDivPrepended($el, "add-button"), "events", NewEventsCanvas);
    NewCanvasAddButton(newNamedDivPrepended($el, "add-button"), "tasks", NewTasksCanvas);
  }

  function NewEventsCanvas() {
    return new Canvas(newNamedDivPrepended($canvases, "canvas"), searchEventsCallback);
  }

  function NewTasksCanvas() {
    var canvas = new TasksCanvas(newNamedDivPrepended($canvases, "canvas"), function() {});
    canvas.Load();
    return canvas;
  }

  function searchEventsCallback(canvas) {
    var canvas2 = NewEventsCanvas();
    canvas2.Search(canvas.SearchCriteria());
    canvas.ResetCriteria();
  }

  setUp();

  return {
    NewEventsCanvas: NewEventsCanvas,
    NewTasksCanvas: NewTasksCanvas,
  };
}

function NewCanvasAddButton($el, title, clickCallback) {
  function setUp() {
    $el.html("<button>+ "+title+"</button>").find("button").click(clickCallback);
  }

  setUp();

  return {};
}

function NewCanvasDeleteButton($el, clickCallback) {
  function setUp() {
    $el.html("<button>-</button>").find("button").click(clickCallback);
  }

  setUp();

  return {};
}

function NewSearchCriteriaExpandingInput($input) {
  function setUp() {
    $input
      .focus(function() { $(this).addClass("expanded"); })
      .blur(function() {
        var $el = $(this);
        setTimeout(function() { $el.removeClass("expanded"); }, 100);
      });
  }

  setUp();

  return {};
}

function NewSearchCriteriaClearButton($input) {
  function setUp() {
    var $button = $("<a class='search-criteria-clear-button'>x</a>").click(function(event) {
      event.preventDefault();
      $input.val("");
      $input.focus();
    });
    $input.after($button);
  }

  setUp();

  return {};
}

function NewMoreEventsButton($el, clickCallback) {
  var $button = null;

  function setUp() {
    $button = $el.html("<button>More events...</button>").find("button");
    $button.click(clickCallback);
    $button.hide(); // default to hide
  }

  function show() { $button.show(); }
  function hide() { $button.hide(); }

  setUp();

  return {
    Show: show,
    Hide: hide,
  };
}

function Canvas($el, searchCallback) {
  var obj = {};
  var currCriteria = new EmptySearchCriteria();

  var $eventsTable = null;
  var moreEventsButton = null;
  var lastEventID = null;

  function setUp() {
    $el.html($("#canvas-tmpl").html());

    $el.find("form").submit(function search(event) {
      event.preventDefault();
      searchCallback(obj);
    });

    $el.find("form input").each(function() {
      NewSearchCriteriaExpandingInput($(this));
      NewSearchCriteriaClearButton($(this));
    });

    $eventsTable = $('<table></table>').appendTo($el);

    $el.on("click", "a[data-query]", function(event) {
      event.preventDefault();
      // todo represent as a object
      var query = $(event.target).data("query");
      var val = $(event.target).data("value");
      $el.find("form").find("input[name='"+query+"']").val(val).focus();
      searchCallback(obj);
    });

    NewHoverableEvents($el);

    NewCanvasDeleteButton(newNamedDivPrepended($el, "delete-button"), function() {
      $el.remove();
    });

    moreEventsButton = NewMoreEventsButton(
      newNamedDiv($el, "more-events-button"), obj.LoadMoreEvents);
  }

  obj.SearchCriteria = function() {
    return new SearchCriteria($el.find("form"));
  };

  obj.Search = function(criteria) {
    criteria.ApplyToForm($el.find("form"));
    criteria.ApplyFocusToForm($el.find("form"));
    currCriteria = criteria;
    $.post("/api/events", criteria.AsQuery()).done(setEvents).fail(showError);
  };

  obj.ResetCriteria = function() {
    currCriteria.ApplyToForm($el.find("form"));
  };

  obj.LoadMoreEvents = function() {
    var criteria = obj.SearchCriteria(); // todo clone currCriteria?
    criteria.SetBeforeID(lastEventID);
    $.post("/api/events", criteria.AsQuery()).done(addEvents).fail(showError);
  };

  function setEvents(data) {
    if (data.Tables[0].Rows.length == 0) {
      $eventsTable.html($("#no-events-tmpl").html());
    } else {
      addEvents(data);
    }
  }

  function addEvents(data) {
    if (data.Tables[0].Rows.length == 0) {
      moreEventsButton.Hide();
    } else {
      var eventsHtml = '';

      data.Tables[0].Rows.forEach(function(apiEvent) {
        eventsHtml += buildEventTmpl(apiEvent);
        lastEventID = apiEvent.id;
      });

      $eventsTable.append(eventsHtml);
      showHideMoreEventsButton();
    }
  }

  function showError() {
    moreEventsButton.Hide();
    $eventsTable.append($("#error-events-tmpl").html());
  }

  function showHideMoreEventsButton() {
    var criteria = obj.SearchCriteria(); // todo clone currCriteria?
    criteria.SetBeforeID(lastEventID);

    $.post("/api/events", criteria.AsQuery()).done(function(data) {
      if (data.Tables[0].Rows.length == 0) {
        moreEventsButton.Hide();
      } else {
        moreEventsButton.Show();
      }
    }).fail(showError);
  }

  var eventHtml = $('#event-tmpl').html();
  var eventKeys = ["action", "context", "deployment", "error", "id",
    "instance", "object_name", "object_type", "task_id", "time", "user"];

  function buildEventTmpl(apiEvent) {
    var eventHtml2 = eventHtml;
    eventKeys.forEach(function(key) {
      eventHtml2 = eventHtml2.replace(new RegExp('{' + key + '}', 'g'), apiEvent[key]);
    });
    return eventHtml2;
  }

  setUp();

  return obj;
}

function NewHoverableEvents($el) {
  var $selected = $el;
  var className = "hover";

  $el.on("mouseover", "table tr", function(event) {
    var $tr = $(event.target).parent("tr");
    if ($tr.length == 0) return;

    $selected.removeClass(className);
    $selected = $tr;

    // Example: "200 <- 199" hovering over 200
    var ids = String($tr.data("event-id")).split(" <- ");
    if (ids.length == 2) {
      $selected = $selected.add($el.find("tr[data-event-id='"+ids[1]+"']"));
    }

    $selected.addClass(className);
  });

  return {};
}

function EmptySearchCriteria() {
  return {
    AsQuery: function() { return "" },
    ApplyToForm: function($el) { $el[0].reset(); },
    ApplyFocusToForm: function($el) {},
  }
}

function SearchCriteria($el) {
  var data = {};
  var focusedInputName = null;

  var keys = ["after", "before", "event-user", "action",
    "object-type", "object-name", "task", "instance", "deployment"];

  function setUp() {
    keys.forEach(function(key) {
      data[key] = $el.find("input[name='"+key+"']").val();
    });

    var $focused = $el.find("input:focus");
    if ($focused.length > 0) {
      focusedInputName = $focused.attr("name");
    }
  }

  function AsQuery() { return data; }

  function ApplyToForm($el2) {
    Object.keys(data).forEach(function(key) {
      $el2.find("input[name='"+key+"']").val(data[key]);
    });
  }

  function ApplyFocusToForm($el2) {
    if (focusedInputName) {
      $el2.find("input[name='"+focusedInputName+"']").focus();
    }
  }

  function SetBeforeID(id) {
    data["before-id"] = id;
  }

  setUp();

  return {
    AsQuery: AsQuery,
    ApplyToForm: ApplyToForm,
    ApplyFocusToForm: ApplyFocusToForm,
    SetBeforeID: SetBeforeID,
  }
}

function TasksCanvas($el, callback_todo) {
  var obj = {};

  function setUp() {
    $el.html($("#tasks-tmpl").html());

    $table = $('<table></table>').appendTo($el);

    NewCanvasDeleteButton(newNamedDivPrepended($el, "delete-button"), function() {
      $el.remove();
    });
  }

  obj.Load = function() {
    $.post("/api/tasks", {"recent": "200", "all": null}).done(setTasks).fail(showError);
  };

  function setTasks(data) {
    if (data.Tables[0].Rows.length > 0) {
      var html = '';

      data.Tables[0].Rows.forEach(function(apiTask) {
        html += buildEventTmpl(apiTask);
        lastEventID = apiTask.id;
      });

      $table.append(html);
    }
  }

  function showError() {
    $table.append($("#error-tasks-tmpl").html());
  }

  var taskHtml = $('#task-tmpl').html();
  var taskKeys = ["deployment", "description", "id",
    "last_activity_at", "result", "started_at", "state", "user"];

  function buildEventTmpl(taskEvent) {
    var taskHtml2 = taskHtml;
    taskKeys.forEach(function(key) {
      taskHtml2 = taskHtml2.replace(new RegExp('{' + key + '}', 'g'), taskEvent[key]);
    });
    return taskHtml2;
  }

  setUp();

  return obj;
}

function newNamedDiv($el, className) {
  return $el.append("<div class='"+className+"'></div>").find("div:last")
}

function newNamedDivPrepended($el, className) {
  return $el.prepend("<div class='"+className+"'></div>").find("div:first")
}

function main() {
  var collection = new CanvasCollection(newNamedDiv($("#root"), "canvas-collection"));

  // start by default with new canvas with all results
  var firstCanvas = collection.NewEventsCanvas();
  firstCanvas.Search(new EmptySearchCriteria());
}

window.addEventListener("load", function load(event){
  window.removeEventListener("load", load, false);
  main();
}, false);

</script>

<style>
.tmpl { display: none; }
button { cursor: pointer; }
.add-button, form, .canvas table, .more-events-button { margin-bottom: 10px; }
input[type="text"], button { font-size: 18px; }
input::placeholder { color: #ccc; }
.canvas input { width: 100px; }
.canvas input.expanded { width: 300px; }
.delete-button { float: right; }
.canvas {
  padding-top: 10px;
  border-top: 2px solid #efefef;
}
table {
  border-spacing: 0;
  border-collapse: collapse;
}
td {
  border: 1px solid #f1f1f1;
  vertical-align: top;
  padding: 0 5px;
}
td.time { width: 230px; }
td.context, td.error { width: 30px; }
td.context span,
td.error span {
  display: inline-block;
  width: 20px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.search-criteria-clear-button {
  background: none;
  border: none;
  vertical-align: top;
  padding: 5px;
  font-size: 12px;
  font-family: system-ui;
  cursor: pointer;
}
table tr.hover { background: yellow; }
</style>

</body>
</html>
`
