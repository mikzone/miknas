<template>
  <div></div>
</template>
<script>
export default {
  name: 'MdcAceEditor',
  props: {
    modelValue: {
      type: String,
      required: true,
    },
    aceLang: {
      type: String,
      default: 'ace/mode/text',
    },
    theme: {
      type: String,
      default: 'monokai',
    },
    autoScrollToEnd: {
      type: Boolean,
      default: false,
    },
    options:{
      type: Object,
      default: ()=>{return {}},
    },
  },
  emits: ['update:modelValue'],
  data: function () {
    return {
      contentBackup: '',
      prevWidth: 0,
    };
  },
  computed: {},
  watch: {
    modelValue: function (val) {
      if (this.contentBackup !== val) {
        this.editor.session.setValue(val);
        this.contentBackup = val;
        this.tryDelayScrollToEnd();
      }
    },
    theme: function (newTheme) {
      this.editor.setTheme('ace/theme/' + newTheme);
    },
    aceLang: function (newLang) {
      // console.log("newLang", newLang);
      this.editor
        .getSession()
        .setMode(typeof newLang === 'string' ? newLang : newLang);
    },
    options: function (newOption) {
      this.editor.setOptions(newOption);
    },
  },
  beforeCreate: function () {
    // console.log("beforeCreate");
    this.editor = null;
  },
  beforeUnmount: function () {
    if (this.editor) {
      this.editor.destroy();
      this.editor.container.remove();
    }
    // console.log("destroy ace editor");
  },
  mounted: function () {
    var lang = this.aceLang;
    var theme = this.theme || 'monokai';

    var editor = (this.editor = window.ace.edit(this.$el));
    editor.$blockScrolling = Infinity;

    // console.log("CmdEditor mounted", this.$el, editor);

    editor
      .getSession()
      .setMode(typeof lang === 'string' ? lang : lang);
    editor.setTheme('ace/theme/' + theme);
    editor.setOptions({
      autoScrollEditorIntoView: true,
      fontSize: 14,
      showGutter: false,
      enableBasicAutocompletion: true,
      enableLiveAutocompletion: true,
      showPrintMargin: false,
      highlightActiveLine: false,
    });
    if (this.modelValue) editor.setValue(this.modelValue, 1);
    this.tryDelayScrollToEnd();
    this.contentBackup = this.modelValue;

    editor.on('change', () => {
      var content = editor.getValue();
      this.$emit('update:modelValue', content);
      this.contentBackup = content;
    });
    if (this.options) editor.setOptions(this.options);
  },
  methods: {
    tryDelayScrollToEnd: function () {
      if (this.autoScrollToEnd) {
        // setTimeout(() => {
        let endLineNum = this.editor.getSession().getLength();
        this.editor.scrollToLine(endLineNum);
        // }, 1000)
      }
    },
  },
};
</script>
