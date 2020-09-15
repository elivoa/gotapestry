
function getStyleObj(dom) {
    return dom ? window.getComputedStyle(dom) : {};
}

function toHyphenCase(name) {
    return name.replace(/[A-Z]/g, function (char) { return "-" + char.toLowerCase(); });
}

function stringifyStyle(styleObj) {
    return Object.keys(styleObj)
        .reduce(function (pre, styleName) {
            var style = styleObj[styleName];
            return style !== undefined && style !== null && !/^\d+$/.test(styleName)
                ? "" + pre + toHyphenCase(styleName) + ":" + style + ";"
                : pre;
        }, '')
        .replace(/"/g, "'");
}

function stringifyAttrs(el, extraAttrs) {
    var attrsMap = Array.prototype.map
        .call(el.attributes, function (at) {
            return ({
                name: at.nodeName,
                value: at.nodeValue
            });
        })
        .filter(function (it) { return it.value && (!extraAttrs || !extraAttrs[it.name]); });
    var extraNames = (extraAttrs
        ? Object.keys(extraAttrs).map(function (name) { return ({ name: name, value: extraAttrs[name] }); })
        : []);
    return attrsMap.concat(extraNames).reduce(function (pre, _a) {
        var name = _a.name, value = _a.value;
        return value ? "" + pre + name + "=\"" + value + "\" " : pre;
    }, '');
}

function stringifyNode(el) {
    if (el.nodeName === '#text' || el.nodeName === '#comment')
        return el.nodeValue || '';
    var styleStr = stringifyStyle(getStyleObj(el));
    var attrs = stringifyAttrs(el, { style: styleStr });
    var tagName = el.nodeName.toLowerCase();
    if (['br', 'hr', 'input', 'img'].indexOf(tagName) !== -1) {
        return "<" + tagName + " " + attrs + "/>";
    }
    var children = Array.prototype.reduce.call(el.childNodes, function (pre, node) { return pre + stringifyNode(node); }, '');
    return "<" + tagName + " " + attrs + ">" + children + "</" + tagName + ">";
}

/**
 * @param { string | Element} [content]     the content you want print.
 *
 *                                          如果值类型为元素节点，将打印对应节点
 *                                          if content is an Element, it will print this element;
 *
 *                                          如果值类型为字符串，将打印解析得到的 html
 *                                          else if content is a string, it will print this string as a html file.
 *
 * @param { CallPrinterOptions } [options]
 * */
function callPrinter(content, options) {
    if (!content)
        window.print();
    else {
        var iframe_1 = document.createElement('iframe');
        iframe_1.setAttribute('style', 'display: none');
        document.body.appendChild(iframe_1);
        var subWindow_1 = iframe_1.contentWindow;
        var $content = void 0;
        if (typeof content === 'string')
            $content = content;
        else
            $content = stringifyNode(content);
        subWindow_1.document.body.innerHTML = $content;
        setTimeout(function () {
            subWindow_1.print();
            document.body.removeChild(iframe_1);
        }, (options && options.delay) || 100);
    }
}


// real print
function printOrderCommon() {
    console.log("print pages commonly!!!")
    var printarea = $('.order_print')
    callPrinter(printarea[0])
}
