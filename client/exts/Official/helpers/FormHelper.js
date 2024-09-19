// 自动生成表单相关

import FormContainerDlg from '../components/forms/FormContainerDlg.vue';
import MdcTextInput from '../components/forms/types/MdcTextInput.vue';
import MdcNumberInput from '../components/forms/types/MdcNumberInput.vue';
import MdcSelect from '../components/forms/types/MdcSelect.vue';
import MdcTextAutoComplete from '../components/forms/types/MdcTextAutoComplete.vue';
import MdcMarkdown from '../components/forms/types/MdcMarkdown.vue';
import { MikCall } from 'miknas/utils';

export const FormTypes = {
  MdcTextInput,
  MdcNumberInput,
  MdcSelect,
  MdcTextAutoComplete,
  MdcMarkdown,
}

export class DataRule {
  static emptyOr(rules) {
    return function (val) {
      let suc = DataRule.isString(val);
      if (suc !== true) return '必须是字符串';
      if (val.length <= 0) return true;
      for (let ruleFunc of rules) {
        let ret = ruleFunc(val);
        if (ret !== true) return ret;
      }
      return true;
    }
  }
  static isInt(val) {
    return Number.isInteger(val) || '必须是一个整数';
  }
  static isString(val) {
    return typeof val == 'string' || '必须是字符串';
  }
  static notEmpty(val) {
    return !!val || '不能为空';
  }
  static isNotEmptyString(val) {
    let suc = DataRule.isString(val) === true && DataRule.notEmpty(val);
    return suc || '必须是非空字符串';
  }
  static minLength(needVal) {
    return function (val) {
      let check1 = DataRule.isString(val);
      if (check1 !== true) return check1;
      if (val.length < needVal) {
        return `长度不能小于${needVal}`;
      }
      return true;
    }
  }
  static maxLength(needVal) {
    return function (val) {
      let check1 = DataRule.isString(val);
      if (check1 !== true) return check1;
      if (val.length > needVal) {
        return `长度不能大于${needVal}`;
      }
      return true;
    }
  }
  static isJsonTxtOrEmpty(val) {
    let suc = DataRule.isString(val);
    if (suc !== true) return '必须是字符串';
    if (val.length <= 0) return true;
    try {
      JSON.parse(val);
    } catch (err) {
      return `不是合法的json格式: ${err}`;
    }
    return true;
  }
  static isJsonDictOrEmpty(val) {
    let suc = DataRule.isString(val);
    if (suc !== true) return '必须是字符串';
    if (val.length <= 0) return true;
    let FinalVal;
    try {
      FinalVal = JSON.parse(val);
    } catch (err) {
      return `不是合法的json格式: ${err}`;
    }
    if (typeof FinalVal != 'object') return '必须是object类型的json格式';
    return true;
  }
  static isUrl(val) {
    let suc = DataRule.isString(val);
    if (suc !== true) return '必须是字符串';
    if (val.length <= 0) return `不能为空`;
    try {
      new URL(val);
    } catch (err) {
      return `不是合法的URL`;
    }
    return true;
  }
}

export class ConfsHelper {
  constructor(settingConfs) {
    this.settingConfs = settingConfs;
  }

  getSettingConfs() {
    return this.settingConfs;
  }

  getDefaultObj() {
    let ret = {};
    for (let confItem of this.settingConfs) {
      ret[confItem.id] = confItem.default;
    }
    return ret;
  }

  checkValueValid(val, confItem) {
    if (confItem.selectOptions) {
      if (['mult-select'].includes(confItem.element_type)) {
        if (!Array.isArray(val)) return `配置项${confItem.id}的值必须为数组`;
      } else {
        for (let option of confItem.selectOptions) {
          if (option.value == val) return true;
        }
        return `配置项${confItem.id}中没有${val}这个选项可选`;
      }
    }
    if (confItem.rules) {
      for (let ruleFunc of confItem.rules) {
        let ruleResult = ruleFunc(val);
        if (ruleResult !== true) return ruleResult;
      }
    }
    return true;
  }

  filterFromExistData(existObj) {
    // 从existObj上提取成新的配置项
    let ret = {};
    for (let confItem of this.settingConfs) {
      let key = confItem.id;
      if (key in existObj) {
        let prevVal = existObj[key];
        let checkResult = this.checkValueValid(prevVal, confItem);
        if (checkResult === true) {
          ret[key] = prevVal;
          continue;
        }
      }
      ret[key] = confItem.default;
    }
    return ret;
  }

  validateFormdata(existObj) {
    // 校验配置
    for (let confItem of this.settingConfs) {
      let key = confItem.id;
      if (key in existObj) {
        let prevVal = existObj[key];
        let checkResult = this.checkValueValid(prevVal, confItem);
        if (checkResult !== true) return checkResult;
      }
    }
    return true;
  }
}

export async function coOpenFormDlg(props) {
  return await MikCall.coCreateDialog({
    component: FormContainerDlg,
    componentProps: props,
  })
}
