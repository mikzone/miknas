<template>
  <q-page padding class="row">
    <q-space />
    <q-form
      :action="postUrl"
      method="post"
      class="q-gutter-md"
      style="width: 600px; max-width: 90%"
      @submit.prevent="onSubmit"
    >
      <q-input
        v-model="uid"
        filled
        label="用户名 *"
        name="uid"
        lazy-rules
        :rules="[(val) => (val && val.length >= 3) || '用户名长度不足']"
      />

      <q-input
        v-model="name"
        filled
        label="昵称 *"
        name="name"
        lazy-rules
        :rules="[(val) => (val && val.length >= 3) || '用户名长度不足']"
      />

      <q-input
        v-model="oriPwd"
        filled
        type="password"
        label="密码 *"
        lazy-rules
        :rules="[(val) => (val && val.length >= 3) || '密码长度不足']"
      />

      <q-input
        v-model="secPwd"
        filled
        type="password"
        label="再次输入密码 *"
        lazy-rules
        :rules="[
          (val) => (val && val.length >= 3) || '密码长度不足',
          (val) => (val && val == oriPwd) || '两次输入的密码要一致',
        ]"
      />

      <q-input v-show="false" v-model="pwd" type="password" name="pwd" />

      <div>
        <q-btn class="full-width" label="注册" type="submit" color="primary" />
      </div>
    </q-form>
    <q-space />
    <q-inner-loading :showing="isLoading">
      <q-spinner-gears size="50px" color="primary" />
    </q-inner-loading>
  </q-page>
</template>

<script>
import { MikCall } from 'miknas/utils';
import { useExtension } from '../extMain';
import { rsaEncrypt } from '../helpers';
let extsObj = useExtension();
export default {
  data: function () {
    return {
      uid: '',
      name: '',
      oriPwd: '',
      secPwd: '',
      pwd: '',
      postUrl: extsObj.serverUrl('register'),
      isLoading: false,
    };
  },
  computed: {},
  methods: {
    async onSubmit(evt) {
      if (this.oriPwd != this.secPwd) {
        return;
      }

      this.isLoading = true;
      let iRet = await extsObj.mcpost('querySecretToken', {});
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        this.isLoading = false;
        return;
      }
      let ret = iRet.ret;
      let token = ret.token;
      this.pwd = rsaEncrypt(token, this.oriPwd);
      this.$nextTick(() => {
        evt.target.submit();
      });
      this.isLoading = false;
    },
  },
};
</script>
