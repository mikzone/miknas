<template>
  <q-select
    use-input
    hide-selected
    bg-color="white"
    :options="selectOptions"
    input-debounce="0"
    fill-input
    :model-value="modelValue"
    v-bind="$attrs"
    @filter="filterFn"
    @input-value="setModel"
  >
    <template #append>
      <slot name="append"></slot>
    </template>
  </q-select>
</template>
<script>
import { gutil } from 'miknas/utils';
export default {
  name: 'MdcRecentInput',
  inheritAttrs: false,
  props: {
    storeKey: {
      type: String,
      require: true,
    },
    modelValue: String,
  },
  data: function () {
    return {
      allOptions: [],
      selectOptions: [],
    };
  },
  computed: {
    realStoreKey() {
      if (!this.storeKey) {
        throw 'storeKey can not be empty';
      }
      return `RI/${this.storeKey}`;
    },
  },
  mounted: function () {
    this.loadRecentInputValues();
  },
  methods: {
    filterFn: function (val, update, abort) {
      update(() => {
        const needle = val.toLocaleLowerCase();
        this.selectOptions = this.allOptions.filter(
          (v) => v.toLocaleLowerCase().indexOf(needle) > -1
        );
      });
    },
    setModel: function (val) {
      this.$emit('update:modelValue', val);
    },
    loadRecentInputValues() {
      let ret = gutil.getStoreItem(this.realStoreKey, []);
      this.selectOptions = ret;
      this.allOptions = ret;
    },
    addRecentInputValue(newVal) {
      if (!newVal) return;
      let preValues = gutil.getStoreItem(this.realStoreKey, []);
      let ret = [newVal];
      for (let v of preValues) {
        if (v == newVal) continue;
        ret.push(v);
        if (ret.length >= 10) break;
      }
      gutil.setStoreItem(this.realStoreKey, ret);
      this.selectOptions = ret;
      this.allOptions = ret;
    },
  },
};
</script>
