(() => {
  // output/Main/foreign.js
  function alertMsg(msg) {
    alert(msg);
  }

  // output/Main/index.js
  var main = /* @__PURE__ */ alertMsg("hello");

  // <stdin>
  main();
})();
