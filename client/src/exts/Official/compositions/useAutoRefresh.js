import { onBeforeUnmount, ref } from 'vue';

function useAutoRefresh(intvMs, refreshCb, needProgress) {
    const mgr = {
        timer: null,
        refreshTs: ref(null),
        progressTimer: null,
        percentage: ref(0),
        nextTs: ref(null),
        start() {
            if (mgr.timer) return;
            refreshCb();
            let now = Date.now();
            mgr.refreshTs.value = now;
            mgr.nextTs.value = now + intvMs;
            mgr.timer = setTimeout(() => {
                mgr.timer = null;
                mgr.start();
            }, intvMs);
            mgr.tryStartProgressTimer();
        },
        stop() {
            if (mgr.timer) {
                clearTimeout(this.timer);
                mgr.timer = null;
                mgr.nextTs.value = null;
                mgr.percentage.value = 0;
            }
            mgr.tryStopProgressTimer();
        },
        tryStartProgressTimer() {
            if (needProgress && mgr.nextTs.value) {
                mgr.progressTimer = setInterval(() => {
                    let now = Date.now();
                    let nextTs = mgr.nextTs.value;
                    if (nextTs) {
                        mgr.percentage.value = (1 - (nextTs - now) / intvMs) * 100;
                    }
                }, 500);
            }
        },
        tryStopProgressTimer() {
            if (mgr.progressTimer) {
                clearInterval(mgr.progressTimer);
                mgr.progressTimer = null;
            }
        },
        forceRefresh() {
            let oldTimer = mgr.timer;
            if (oldTimer) {
                mgr.stop();
                mgr.start();
            } else {
                mgr.refreshTs.value = Date.now();
                refreshCb();
            }
        },
    }

    onBeforeUnmount(() => {
        mgr.stop();
    });
    return mgr
}

export default useAutoRefresh
