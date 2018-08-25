var currentUri = window.location.pathname;

// 返回/前进浏览器事件
window.onpopstate = function() {
    if (currentUri != window.location.pathname) {
        loadMarkdown(window.location.pathname, false)
    }
};

// 打开节点,参数为li节点的jquery对象
function openNode(node) {
    if (node.find("ul").eq(0).css('display') == "none") {
        node.find("ul").eq(0).show();
        node.find("i").eq(0).attr("class", "am-icon-caret-down");
    }
}

// 替换a标签及img标签的地址，加上/前缀修改为绝对路径
function replaceHrefAndSrc() {
    // 修改a/img标签链接，给相对路径统一加上前缀
    $(document).find("a").each(function(){
        var href = $(this).attr("href");
        if (typeof href != "undefined" && href.length > 0) {
            if (href.substr(0, 1) != "/" && href.substr(0, 1) != "#" && href.substr(0, 4) != "http") {
                href = "/" + href
                if (href.indexOf(".md")) {
                    href = href.replace(".md", "");
                    $(this).attr("href", "javascript:loadMarkdown('" + href + "', true);");
                } else {
                    $(this).attr("href", href);
                }
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
                $(this).attr("src", "/" + src);
            }
        }
    });
}

// 取消当前的li高亮
function cancelAllHighlight() {
    $("#side-markdown-view").find("li").removeClass("active");
}

// 高亮并展开当前打开的地址
function highlightLiByUri(uri) {
    var seli = $("a[href='"+ uri +"']").parent("li");
    seli.addClass("active");
    openNode(seli.parent("ul").parent("li"));
}

// 重新解析markdown内容
function reloadMainMarkdown() {
    var content = $("#main-markdown-content").text()
    if (content.length > 0) {
        $('#main-markdown-view').html(marked($("#main-markdown-content").text()));
        $('#main-markdown-view pre code').each(function(i, block) {
            Prism.highlightElement(block);
        });
        // 生成TOC菜单
        $('#main-markdown-toc').html("");
        new Toc('main-markdown-view', {
            'level'   : 3,
            'class'   : 'toc',
            'targetId': 'main-markdown-toc'
        } );
        if ($('#main-markdown-toc').html().length > 0) {
            var html = $("#main-markdown-view").html().replace("<p>[TOC]</p>", $('#main-markdown-toc').html())
            $("#main-markdown-view").html(html)
        }
    }
    replaceHrefAndSrc();
    updateHelpUrl(window.location.pathname);
}

// 更新文档markdown链接地址
function updateHelpUrl(uri) {
    $("#help-icon").attr("href", "https://gitee.com/johng/gf-doc/tree/master" + uri + ".md");
}

// 请求markdown内容
function loadMarkdown(uri, addState) {
    currentUri = uri;
    cancelAllHighlight();
    // 添加历史记录
    if (addState) {
        window.history.pushState({
            title : document.title,
            uri   : window.location.pathname
        }, document.title, window.location.origin + uri);
    }
    highlightLiByUri(uri);
    // 修改当前标题
    var title = $("a[href='"+ uri +"']").text();
    if (title.length > 0) {
        title += " - " + baseTitle
    } else {
        title += "404 NOT FOUND - " + baseTitle
    }
    document.title = title;
    $("#main-markdown-view").html("<div class=\"loading-small\"></div> Loading...");
    // AJAX读取文档
    $.ajax({
        type     : "get",
        url      : uri,
        dataType : "json",
        success: function(result){
            if (result.code == 1) {
                $("#main-markdown-content").text(result.data);
                reloadMainMarkdown();
            }
        }
    });
}

$(function() {
    reloadMainMarkdown();
    // 修改list样式
    $("#side-markdown-view").find("ul").addClass("am-list am-list-border");
    // 回到顶部
    $('#totop-icon').on('click', function() {
        $("#main-markdown-view").smoothScroll({position: 0, speed: 300});
    });
    $("#main-markdown-view").scroll(function() {
        if ($(this).scrollTop() >= 400) {
          $("#totop-icon").show();
        } else {
          $("#totop-icon").hide();
        }
    });
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
            loadMarkdown($(this).find("a").eq(0).attr('href'), true)
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
    highlightLiByUri(window.location.pathname);

    $("#process-mask").hide();

    // 菜单关闭隐藏
    $("#menu-icon").click(function () {
        if ($("#side-markdown-view").css("display") == "none") {
            $(this).css("left", "340px");
            $("#side-markdown-view").show();
        } else {
            $(this).css("left", "20px");
            $("#side-markdown-view").hide();
        }

    });
});
