import { computed, ref } from 'vue';

export function useMikLoading() {
  const loadingStates = ref({});

  const loadingLabel = computed(() => {
    return Object.keys(loadingStates.value).join(',');
  })

  const isloading = computed(() => {
    return loadingLabel.value.length > 0;
  })


  function addLoadingState(stateName) {
    loadingStates.value[stateName] = true;
  }

  function removeLoadingState(stateName) {
    if (loadingStates.value[stateName])
      delete (loadingStates.value[stateName]);
  }

  return {
    isloading,
    loadingLabel,
    addLoadingState,
    removeLoadingState,
  }
}
