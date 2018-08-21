
window.onpopstate = function() {
    var state = window.history.state;
    console.log(state)
    loadMarkdown(state.uri)
}

// 打开节点,参数为li节点的jquery对象
function openNode(node) {
    if (node.find("ul").eq(0).css('display') == "none") {
        node.find("ul").eq(0).show();
        node.find("i").eq(0).attr("class", "am-icon-caret-down");
    }
}

// 取消当前的li高亮
function cancelAllHighlight() {
    $("#side-markdown-view").find("li").removeClass("active");
}

// 高亮并展开当前打开的地址
function highlightCurrentLi() {
    var seli = $("a[href='"+ window.location.pathname +"']").parent("li")
    seli.addClass("active");
    openNode(seli.parent("ul").parent("li"));
}

// 重新解析markdown内容
function reloadMainMarkdown() {
    $("#main-markdown-view").html("");
    editormd.markdownToHTML("main-markdown-view", {
        markdown        : $("#main-markdown-content").text(),
        htmlDecode      : true,
        toc             : true,
        gfm             : true,
        emoji           : true,
        taskList        : true,
        tex             : true,
        tocDropdown     : false,
        markdownSourceCode : false
    });
}

// 请求markdown内容
function loadMarkdown(uri) {
    cancelAllHighlight()
    // 添加历史记录
    window.history.pushState({
        title : document.title,
        uri   : window.location.pathname
    }, document.title, window.location.origin + uri);
    // 修改当前标题
    document.title = $("a[href='"+ uri +"']").text() + " - " + baseTitle
    $("#main-markdown-view").html("<div class=\"loading-small\"></div> Loading...");
    // AJAX读取文档
    $.ajax({
        type     : "get",
        url      : uri,
        dataType : "json",
        success: function(result){
            if (result.code == 1) {
                $("#main-markdown-content").text(result.data)
                reloadMainMarkdown()
                highlightCurrentLi()
            }
        }
    });
}

$(function() {
    reloadMainMarkdown()
    // 修改a/img标签链接，给相对路径统一加上前缀
    $(document).find("a").each(function(){
        var href = $(this).attr("href");
        if (typeof href != "undefined" && href.length > 0) {
            if (href.substr(0, 1) != "/" && href.substr(0, 1) != "#" && href.substr(0, 4) != "http") {
                $(this).attr("href", "/doc/" + href);
            }
            if (href.substr(0, 4) == "http") {
                $(this).attr("target", "_blank");
            }
        }
    });
    $(document).find("img").each(function(){
        var src = $(this).attr("src");
        if (typeof src != "undefined" && src.length > 0) {
            if (src.substr(0, 1) != "/" && src.substr(0, 4) != "http") {
                $(this).attr("src", "/doc/" + src);
            }
        }
    });

    // 修改list样式
    $("#side-markdown-view").find("ul").addClass("am-list am-list-border");

    // 菜单树形结构处理
    function indentTreeByLi(obj) {
        if ($(obj).find("ul").length > 0) {
            // 添加图标
            $(obj).find("a").eq(0).find("i").attr("class", "am-icon-caret-down");
            // 递归处理树型缩进
            $(obj).find("ul").eq(0).find(">li").each(function () {
                $(this).find("a").each(function () {
                    var v = $(this).css("padding-left");
                    var i = parseInt(v) + 20;
                    $(this).css("padding-left", i + "px");
                });
                if ($(this).find(">ul").length > 0) {
                    indentTreeByLi(this)
                }
            })
        }
    }
    $("#side-markdown-view").find("a").prepend('<i class="empty-i"></i>');
    $("#side-markdown-view").find("ul").eq(0).find(">li").each(function(){
        indentTreeByLi(this)
    });
    // 绑定li点击事件
    $("#side-markdown-view").find("li").click(function(){
        if ($(this).find("ul").length > 0) {
            $(this).find("ul").eq(0).toggle();
            // 改变图标
            if ($(this).find("ul").eq(0).css('display') == "none") {
                $(this).find("i").eq(0).attr("class", "am-icon-caret-right");
            } else {
                $(this).find("i").eq(0).attr("class", "am-icon-caret-down");
            }
        } else {
            loadMarkdown($(this).find("a").eq(0).attr('href'))
        }
        return false
    });
    // 先将所有带子级的li关闭
    $("#side-markdown-view").find("ul").find(">li").each(function () {
        if ($(this).find("ul").length > 0) {
            $(this).trigger("click")
        }
    });
    // 默认将第一级菜单展开
    $("#side-markdown-view").find("ul").eq(0).find(">li").each(function () {
        if ($(this).find("ul").length > 0) {
            $(this).trigger("click")
        }
    });
    // 高亮并展开当前打开的地址
    highlightCurrentLi()

    $("#process-mask").hide();
});
