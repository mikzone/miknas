import { gutil, MikCall } from 'miknas/utils';
import { reactive } from 'vue';
import { ConfsHelper } from '../helpers/FormHelper';

// 表单视图
function useFormView(formConfs, initData) {
  const confsHelpr = new ConfsHelper(formConfs);
  const state = reactive({
    formData: confsHelpr.filterFromExistData(initData),
  });

  const action = {
    tryGetValidFormData(){
      let ret = confsHelpr.validateFormdata(state.formData);
      if (ret === true) return gutil.jsonCopy(state.formData);
      MikCall.sendErrorTips(ret);
    }
  };

  return {
    formConfs,
    state,
    action,
  }
}

export default useFormView
