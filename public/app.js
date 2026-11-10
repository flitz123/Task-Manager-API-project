const API = "/tasks";
let tasks = [];
let viewFilter = "all";
let statusFilter = "all";
let searchText = "";
let editingId = null;

const $ = (selector) => document.querySelector(selector);
const $$ = (selector) => [...document.querySelectorAll(selector)];
const today = () => new Date().toISOString().slice(0, 10);

async function api(path, options = {}) {
  const response = await fetch(path, { headers: {"Content-Type":"application/json"}, ...options });
  if (!response.ok && response.status !== 204) {
    const data = await response.json().catch(() => ({}));
    throw new Error(data.error || "Request failed");
  }
  return response.status === 204 ? null : response.json();
}

async function loadTasks() {
  try { tasks = await api(API); }
  catch (error) {
    const saved = localStorage.getItem("task-manager-fallback");
    tasks = saved ? JSON.parse(saved) : [];
    toast("Offline mode: changes are saved locally");
  }
  render();
}

function saveFallback() {
  localStorage.setItem("task-manager-fallback", JSON.stringify(tasks));
}

function escapeHTML(value = "") {
  return value.replace(/[&<>"']/g, char => ({ "&":"&amp;", "<":"&lt;", ">":"&gt;", '"':"&quot;", "'":"&#39;" }[char]));
}

function render() {
  const total = tasks.length;
  const completed = tasks.filter(t => t.completed).length;
  const active = total - completed;
  const todayTasks = tasks.filter(t => t.dueDate === today() && !t.completed).length;
  $("#all-count").textContent = total;
  $("#today-count").textContent = todayTasks;
  $("#active-count").textContent = active;
  $("#completed-count").textContent = completed;
  $("#stat-total").textContent = total;
  $("#stat-today").textContent = todayTasks;
  $("#stat-progress").textContent = total ? Math.round(completed / total * 100) + "%" : "0%";

  let visible = tasks.filter(task => {
    if (viewFilter === "today") return task.dueDate === today();
    if (viewFilter === "active") return !task.completed;
    if (viewFilter === "completed") return task.completed;
    return true;
  });
  visible = visible.filter(task => statusFilter === "all" || (statusFilter === "completed" ? task.completed : !task.completed));
  if (searchText) visible = visible.filter(task => `${task.title} ${task.description || ""} ${task.category || ""}`.toLowerCase().includes(searchText));

  const sort = $("#sort-select").value;
  visible.sort((a,b) => {
    if (sort === "priority") return {high:0,medium:1,low:2}[a.priority || "medium"] - {high:0,medium:1,low:2}[b.priority || "medium"];
    if (sort === "due") return (a.dueDate || "9999").localeCompare(b.dueDate || "9999");
    if (sort === "oldest") return (a.createdAt || "").localeCompare(b.createdAt || "");
    return (b.createdAt || "").localeCompare(a.createdAt || "");
  });

  const list = $("#task-list");
  list.innerHTML = visible.map(task => {
    const category = task.category ? `<span class="category">${escapeHTML(task.category)}</span>` : "";
    const due = task.dueDate ? `<span class="due ${task.dueDate < today() && !task.completed ? "overdue" : ""}">◷ ${formatDate(task.dueDate)}</span>` : "";
    return `<article class="task-row ${task.completed ? "completed" : ""}" data-id="${task.id}">
      <input class="task-check" type="checkbox" ${task.completed ? "checked" : ""} aria-label="Toggle ${escapeHTML(task.title)}">
      <div class="task-main"><div class="task-title">${escapeHTML(task.title)}</div><div class="task-meta">${category}<span class="priority ${task.priority || "medium"}">● ${(task.priority || "medium")}</span></div></div>
      ${due}<div class="task-actions"><button class="icon-button edit-task" type="button" aria-label="Edit task">✎</button><button class="icon-button delete-task" type="button" aria-label="Delete task">⌫</button></div>
    </article>`;
  }).join("");
  $("#empty-state").hidden = visible.length > 0;
  $("#task-list").hidden = visible.length === 0;
}

function formatDate(value) {
  const date = new Date(value + "T00:00:00");
  return date.toLocaleDateString(undefined, { month:"short", day:"numeric" });
}

function openModal(task = null) {
  editingId = task?.id || null;
  $("#modal-title").textContent = task ? "Edit task" : "Create task";
  $("#save-label").textContent = task ? "Save changes" : "Create task";
  $("#task-title").value = task?.title || "";
  $("#task-description").value = task?.description || "";
  $("#task-priority").value = task?.priority || "medium";
  $("#task-due").value = task?.dueDate || "";
  $("#task-category").value = task?.category || "";
  $("#modal-backdrop").hidden = false;
  $("#task-title").focus();
}

function closeModal() {
  $("#modal-backdrop").hidden = true;
  editingId = null;
  $("#task-form").reset();
}

function toast(message) {
  const el = $("#toast");
  el.textContent = message;
  el.classList.add("show");
  clearTimeout(toast.timer);
  toast.timer = setTimeout(() => el.classList.remove("show"), 2300);
}

$("#task-form").addEventListener("submit", async event => {
  event.preventDefault();
  const payload = {
    title: $("#task-title").value.trim(),
    description: $("#task-description").value.trim(),
    priority: $("#task-priority").value,
    dueDate: $("#task-due").value,
    category: $("#task-category").value.trim(),
    completed: editingId ? tasks.find(t => t.id === editingId)?.completed || false : false
  };
  if (!payload.title) return;
  try {
    if (editingId) await api(API + "/" + editingId, { method:"PUT", body:JSON.stringify(payload) });
    else await api(API, { method:"POST", body:JSON.stringify(payload) });
    toast(editingId ? "Task updated" : "Task created");
  } catch (error) {
    if (editingId) {
      const index = tasks.findIndex(t => t.id === editingId);
      if (index >= 0) tasks[index] = {...tasks[index], ...payload, updatedAt:new Date().toISOString()};
    } else {
      tasks.push({...payload, id:Math.max(0,...tasks.map(t=>t.id))+1, createdAt:new Date().toISOString(), updatedAt:new Date().toISOString()});
    }
    saveFallback();
    toast("Saved locally");
  }
  closeModal();
  await loadTasks();
});

$("#task-list").addEventListener("change", async event => {
  if (!event.target.classList.contains("task-check")) return;
  const row = event.target.closest(".task-row");
  const task = tasks.find(t => t.id === Number(row.dataset.id));
  const updated = {...task, completed:event.target.checked};
  try {
    await api(API + "/" + task.id, { method:"PUT", body:JSON.stringify(updated) });
  } catch {
    Object.assign(task, updated);
    saveFallback();
  }
  await loadTasks();
});

$("#task-list").addEventListener("click", async event => {
  const row = event.target.closest(".task-row");
  if (!row) return;
  const task = tasks.find(t => t.id === Number(row.dataset.id));
  if (event.target.classList.contains("edit-task")) openModal(task);
  if (event.target.classList.contains("delete-task")) {
    if (!confirm(`Delete "${task.title}"?`)) return;
    try { await api(API + "/" + task.id, { method:"DELETE" }); }
    catch { tasks = tasks.filter(t => t.id !== task.id); saveFallback(); }
    toast("Task deleted");
    await loadTasks();
  }
});

$$(".nav-item").forEach(button => button.addEventListener("click", () => {
  viewFilter = button.dataset.filter;
  $$(".nav-item").forEach(item => item.classList.remove("active"));
  button.classList.add("active");
  const titles = {all:"All tasks",today:"Today",active:"Active tasks",completed:"Completed"};
  $("#view-title").textContent = titles[viewFilter];
  $("#view-eyebrow").textContent = viewFilter === "all" ? "Workspace" : "Filtered view";
  render();
  $(".sidebar").classList.remove("open");
}));

$$(".chip").forEach(button => button.addEventListener("click", () => {
  statusFilter = button.dataset.status;
  $$(".chip").forEach(item => item.classList.remove("active"));
  button.classList.add("active");
  render();
}));

$("#search-input").addEventListener("input", event => {
  searchText = event.target.value.toLowerCase();
  render();
});
$("#sort-select").addEventListener("change", render);
$("#add-button").addEventListener("click", () => openModal());
$(".empty-add").addEventListener("click", () => openModal());
$("#close-modal").addEventListener("click", closeModal);
$("#cancel-modal").addEventListener("click", closeModal);
$("#modal-backdrop").addEventListener("click", event => { if (event.target === $("#modal-backdrop")) closeModal(); });
$("#mobile-menu").addEventListener("click", () => $(".sidebar").classList.toggle("open"));

$("#clear-button").addEventListener("click", async () => {
  const completed = tasks.filter(t => t.completed);
  if (!completed.length) return;
  if (!confirm(`Clear ${completed.length} completed task${completed.length === 1 ? "" : "s"}?`)) return;
  await Promise.all(completed.map(task => api(API + "/" + task.id, { method:"DELETE" }).catch(() => null)));
  tasks = tasks.filter(t => !t.completed);
  saveFallback();
  toast("Completed tasks cleared");
  await loadTasks();
});

$("#theme-button").addEventListener("click", () => {
  document.body.classList.toggle("dark");
  localStorage.setItem("task-manager-theme", document.body.classList.contains("dark") ? "dark" : "light");
  $("#theme-icon").textContent = document.body.classList.contains("dark") ? "☀" : "◐";
});

document.addEventListener("keydown", event => {
  if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") {
    event.preventDefault();
    $("#search-input").focus();
  }
  if (event.key === "Escape") closeModal();
});

if (localStorage.getItem("task-manager-theme") === "dark") {
  document.body.classList.add("dark");
  $("#theme-icon").textContent = "☀";
}
loadTasks();
