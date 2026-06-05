(function () {
    'use strict';

    let i18nInitialized = false;
    let paymentMethods = [];
    let selectedCurrency = '';
    let selectedNetwork = '';
    let selectedMethod = null;
    let currentPayment = {};
    let paymentConfig = {};
    let toastTimer = null;

    const TOKEN_LOGOS = {
        USDT: '/payment/assets/img/token-usdt.svg',
        USDC: '/payment/assets/img/token-usdc.svg',
        TRX: '/payment/assets/img/token-trx.svg',
        ETH: '/payment/assets/img/token-eth.svg',
        BNB: '/payment/assets/img/token-bnb.svg'
    };
    const NETWORK_LOGOS = {
        TRON: '/payment/assets/img/network-tron.svg',
        TRC20: '/payment/assets/img/network-tron.svg',
        ETHEREUM: '/payment/assets/img/network-ethereum.svg',
        ERC20: '/payment/assets/img/network-ethereum.svg',
        BSC: '/payment/assets/img/network-bsc.svg',
        BEP20: '/payment/assets/img/network-bsc.svg',
        POLYGON: '/payment/assets/img/network-polygon.svg',
        ARBITRUM: '/payment/assets/img/network-arbitrum.svg',
        'ARBITRUM-ONE': '/payment/assets/img/network-arbitrum.svg',
        BASE: '/payment/assets/img/network-base.svg',
        SOLANA: '/payment/assets/img/network-solana.svg',
        APTOS: '/payment/assets/img/network-aptos.svg',
        'X LAYER': '/payment/assets/img/network-xlayer.svg',
        XLAYER: '/payment/assets/img/network-xlayer.svg',
        'X-LAYER': '/payment/assets/img/network-xlayer.svg',
        PLASMA: '/payment/assets/img/network-plasma.svg'
    };

    let currencySelectEl = null;
    let networkSelectEl = null;
    let payButton = null;
    let amountValue = null;
    let addressValue = null;
    let copyAddressButton = null;
    let qrcodeEl = null;
    let paymentHintEl = null;
    let exchangeRateEl = null;
    let selectedCurrencyEl = null;
    let selectedNetworkEl = null;
    let switchCounterEl = null;

    let countdownTimer = null;
    let statusCheckTimer = null;
    let totalSeconds = 0;
    let tradeId = '';

    function initI18n() {
        if (i18nInitialized) {
            return Promise.resolve();
        }

        let currentLang = localStorage.getItem('payment_language') || navigator.language.split('-')[0] || 'en';
        if (currentLang !== 'zh' && currentLang !== 'en') {
            currentLang = 'en';
        }

        return new Promise(function (resolve) {
            if (typeof i18next === 'undefined') {
                console.error('i18next is not loaded');
                resolve();
                return;
            }

            i18next.init({
                lng: currentLang,
                debug: false,
                resources: {}
            }, function (err) {
                if (err) {
                    console.error('i18next initialization failed:', err);
                    resolve();
                    return;
                }

                fetch('/payment/assets/locales/' + currentLang + '.json')
                    .then(function (r) {
                        return r.json();
                    })
                    .then(function (translations) {
                        i18next.addResourceBundle(currentLang, 'translation', translations);
                        i18nInitialized = true;
                        updateContent();

                        const languageSwitcher = document.getElementById('languageSwitcher');
                        if (languageSwitcher) {
                            languageSwitcher.value = currentLang;
                        }
                        resolve();
                    })
                    .catch(function (error) {
                        console.error('Failed to load language file:', error);
                        resolve();
                    });
            });
        });
    }

    function changeLanguage(lang) {
        if (!i18next) return;

        const hasLanguage = i18next.hasResourceBundle(lang, 'translation');
        if (!hasLanguage) {
            fetch('/payment/assets/locales/' + lang + '.json')
                .then(function (r) {
                    return r.json();
                })
                .then(function (translations) {
                    i18next.addResourceBundle(lang, 'translation', translations);
                    return i18next.changeLanguage(lang);
                })
                .then(function () {
                    localStorage.setItem('payment_language', lang);
                    updateContent();
                    updateUI();
                })
                .catch(function (error) {
                    console.error('Failed to load language file:', error);
                });
            return;
        }

        i18next.changeLanguage(lang, function (err) {
            if (err) {
                console.error('Language change failed:', err);
                return;
            }
            localStorage.setItem('payment_language', lang);
            updateContent();
            updateUI();
        });
    }

    function replacePlaceholders(text) {
        const active = selectedMethod || currentPayment || {};
        const token = (active.currency || currentPayment.currency || paymentConfig.crypto || '--').toString().toUpperCase();
        const network = (active.network || currentPayment.network || '--').toString().toUpperCase();
        const networkName = (active.token_net_name || currentPayment.token_net_name || network).toString().toUpperCase();

        return text
            .replace(/\{\{token\}\}/g, token)
            .replace(/\{\{network\}\}/g, network)
            .replace(/\{\{networkName\}\}/g, networkName)
            .replace(/\{\{warningToken\}\}/g, network);
    }

    function updateContent() {
        if (!i18next) return;

        document.querySelectorAll('[data-i18n]').forEach(function (element) {
            if (element.querySelector('.options-list')) {
                return;
            }

            const key = element.getAttribute('data-i18n');
            if (!key) return;

            if (key.startsWith('[')) {
                const matches = key.match(/\[(.+?)\](.+)/);
                if (!matches) return;

                const attr = matches[1];
                const transKey = matches[2];
                let translation = i18next.t(transKey);
                translation = replacePlaceholders(translation);

                if (attr === 'html') {
                    element.innerHTML = translation;
                } else {
                    element.setAttribute(attr, translation);
                }
                return;
            }

            let translation = i18next.t(key);
            translation = replacePlaceholders(translation);
            if (element.tagName === 'INPUT' || element.tagName === 'TEXTAREA') {
                element.placeholder = translation;
            } else {
                element.innerHTML = translation;
            }
        });
    }

    function t(key, defaultValue) {
        if (i18nInitialized && i18next) {
            return replacePlaceholders(i18next.t(key));
        }
        return defaultValue || key;
    }

    function normalizeLogoKey(value) {
        return (value || '').toString().trim().toUpperCase();
    }

    function logoInitial(value) {
        const key = normalizeLogoKey(value);
        return key ? key.slice(0, 3) : '--';
    }

    function logoMeta(type, value) {
        const key = normalizeLogoKey(value);
        const logos = type === 'network' ? NETWORK_LOGOS : TOKEN_LOGOS;
        return {
            src: logos[key] || '',
            fallback: logoInitial(key)
        };
    }

    function createLogoNode(type, value, className) {
        const meta = logoMeta(type, value);
        const wrap = document.createElement('span');
        wrap.className = className + ' logo-mark';

        if (meta.src) {
            const img = document.createElement('img');
            img.src = meta.src;
            img.alt = normalizeLogoKey(value);
            img.loading = 'lazy';
            wrap.appendChild(img);
            return wrap;
        }

        wrap.classList.add('logo-fallback');
        wrap.textContent = meta.fallback;
        return wrap;
    }

    function logoTypeForSelect(parent) {
        return parent === networkSelectEl ? 'network' : 'token';
    }

    function renderSelectedLabel(parent, value, text) {
        const sp = parent ? parent.querySelector('.select-label') : null;
        if (!sp) return;

        sp.removeAttribute('data-i18n');
        sp.textContent = '';
        sp.appendChild(createLogoNode(logoTypeForSelect(parent), value, 'selected-logo'));

        const labelText = document.createElement('span');
        labelText.className = 'selected-text';
        labelText.textContent = text || value || '--';
        sp.appendChild(labelText);
    }

    function setupDropdown(el) {
        if (!el) return;

        el.addEventListener('click', function (e) {
            if (el.classList.contains('locked')) return;
            e.stopPropagation();
            document.querySelectorAll('.coin-select').forEach(function (other) {
                if (other !== el) other.classList.remove('active');
            });
            el.classList.toggle('active');
        });
    }

    function initSelection() {
        if (!paymentMethods || paymentMethods.length === 0) return;

        const currencies = Array.from(new Set(paymentMethods.map(function (m) {
            return m.currency;
        })));

        renderOptions(currencySelectEl, currencies.map(function (c) {
            return {value: c, label: c, badge: c === 'USDT' || c === 'USDC' ? t('payment.preferred', 'Preferred') : ''};
        }), function (val) {
            selectedCurrency = val;
            selectedNetwork = '';
            selectedMethod = null;
            updateNetworkOptions();
        });

        const preferred = findCurrentMethod();
        if (preferred && hasConcretePayment(currentPayment)) {
            selectedCurrency = preferred.currency;
            selectOption(currencySelectEl, selectedCurrency);
            updateNetworkOptions(preferred);
            return;
        }

        selectedCurrency = '';
        selectedNetwork = '';
        selectedMethod = null;
        resetSelectLabel(currencySelectEl, 'payment.selectCurrency');
        resetSelectLabel(networkSelectEl, 'payment.selectNetwork');
        renderOptions(networkSelectEl, [], null);
        updateUI();
        updateContent();
    }

    function findCurrentMethod() {
        if (!currentPayment || !currentPayment.currency) return null;

        return paymentMethods.find(function (m) {
            return m.currency === currentPayment.currency &&
                (m.network === currentPayment.network || m.token_net_name === currentPayment.token_net_name);
        }) || null;
    }

    function updateNetworkOptions(preferredMethod) {
        const methods = paymentMethods.filter(function (m) {
            return m.currency === selectedCurrency;
        });

        renderOptions(networkSelectEl, methods.map(function (m) {
            return {
                value: m.token_net_name,
                label: m.token_net_name.toUpperCase(),
                fullData: m
            };
        }), function (val, item) {
            selectedNetwork = val;
            selectedMethod = item.fullData;
            updateUI();
            updateContent();
        });

        const method = preferredMethod && preferredMethod.currency === selectedCurrency ? preferredMethod : null;
        if (method) {
            selectedNetwork = method.token_net_name;
            selectedMethod = method;
            selectOption(networkSelectEl, selectedNetwork);
        } else {
            selectedNetwork = '';
            selectedMethod = null;
            resetSelectLabel(networkSelectEl, 'payment.selectNetwork');
        }

        updateUI();
        updateContent();
    }

    function renderOptions(parent, items, onSelect) {
        if (!parent) return;

        let list = parent.querySelector('.options-list');
        if (!list) {
            list = document.createElement('div');
            list.className = 'options-list';
            parent.appendChild(list);
        }

        list.innerHTML = '';
        items.forEach(function (item) {
            const itemEl = document.createElement('div');
            itemEl.className = 'option-item';
            itemEl.setAttribute('data-value', item.value);

            const optionMain = document.createElement('span');
            optionMain.className = 'option-main';
            optionMain.appendChild(createLogoNode(logoTypeForSelect(parent), item.value, 'option-logo'));

            const optionText = document.createElement('span');
            optionText.className = 'option-text';
            optionText.textContent = item.label;
            optionMain.appendChild(optionText);
            itemEl.appendChild(optionMain);

            if (item.badge) {
                const badge = document.createElement('span');
                badge.className = 'option-badge';
                badge.textContent = item.badge;
                itemEl.appendChild(badge);
            }

            itemEl.addEventListener('click', function (e) {
                e.stopPropagation();
                parent.classList.remove('active');
                selectOption(parent, item.value);
                if (onSelect) onSelect(item.value, item);
            });

            list.appendChild(itemEl);
        });
    }

    function selectOption(parent, value) {
        if (!parent) return;

        const items = parent.querySelectorAll('.option-item');
        let selected = null;
        items.forEach(function (item) {
            item.classList.remove('selected');
            if (item.getAttribute('data-value') === value) selected = item;
        });

        if (!selected) return;

        selected.classList.add('selected');
        const textEl = selected.querySelector('.option-text');
        const text = textEl ? textEl.textContent : value;
        renderSelectedLabel(parent, value, text || value);
    }

    function resetSelectLabel(parent, key) {
        if (!parent) return;
        const sp = parent.querySelector('.select-label');
        if (sp) {
            sp.setAttribute('data-i18n', key);
            sp.textContent = t(key, key);
        }
    }

    function hasConcretePayment(payment) {
        return !!(payment && payment.address);
    }

    function isInitialSelectionMode() {
        return !hasConcretePayment(currentPayment);
    }

    function selectedMatchesCurrent() {
        if (!selectedMethod || !hasConcretePayment(currentPayment)) return false;
        return selectedMethod.currency === currentPayment.currency &&
            (selectedMethod.network === currentPayment.network || selectedMethod.token_net_name === currentPayment.token_net_name);
    }

    function remainingSwitchCount() {
        const value = parseInt(currentPayment.remaining_switch_count, 10);
        return isNaN(value) ? 0 : value;
    }

    function maxSwitchCount() {
        const value = parseInt(currentPayment.max_switch_count, 10);
        return isNaN(value) ? 0 : value;
    }

    function isWaitingOrder() {
        const status = parseInt(currentPayment.status, 10);
        return !status || status === 1;
    }

    function isSwitchLocked() {
        return !isWaitingOrder() || (!!currentPayment.address && remainingSwitchCount() <= 0);
    }

    function updateUI() {
        const hasReadyPayment = selectedMatchesCurrent();
        const initialSelection = isInitialSelectionMode();
        const active = hasReadyPayment ? currentPayment : (selectedMethod || currentPayment || {});
        const hasSelected = !!selectedMethod;
        const cryptoIcon = document.querySelector('.crypto-icon');
        const container = document.querySelector('.cashier-container');

        if (container) {
            container.classList.toggle('is-initial-selection', initialSelection);
            container.classList.toggle('is-method-selected', hasSelected);
            container.classList.toggle('is-payment-ready', hasConcretePayment(currentPayment));
        }

        if (cryptoIcon) {
            cryptoIcon.textContent = '';
            cryptoIcon.appendChild(createLogoNode('token', active.currency || '', 'hero-logo'));
        }
        if (amountValue) {
            const amount = active.actual_amount && active.actual_amount !== '0' ? active.actual_amount : '--';
            const currency = active.currency || '';
            amountValue.textContent = currency ? amount + ' ' + currency : '--';
        }
        if (selectedCurrencyEl) {
            selectedCurrencyEl.textContent = active.currency || '--';
        }
        if (selectedNetworkEl) {
            selectedNetworkEl.textContent = (active.token_net_name || active.network || '--').toString().toUpperCase();
        }
        if (exchangeRateEl) {
            exchangeRateEl.textContent = active.exchange_rate || '--';
        }
        if (switchCounterEl) {
            switchCounterEl.textContent = remainingSwitchCount() + ' / ' + maxSwitchCount();
        }

        renderAddressAndQr(hasReadyPayment);
        updatePayButton(hasSelected);
        updateSelectorLockState();
    }

    function renderAddressAndQr(hasReadyPayment) {
        if (hasReadyPayment) {
            if (addressValue) addressValue.textContent = currentPayment.address;
            if (copyAddressButton) copyAddressButton.disabled = false;
            renderQrCode(currentPayment.address);
            if (paymentHintEl) paymentHintEl.textContent = t('payment.paymentReady', 'Payment details are ready. Transfer the exact amount below.');
            return;
        }

        const pendingAddressKey = isInitialSelectionMode() ? 'payment.addressPendingStart' : 'payment.addressPending';
        const pendingHintKey = isInitialSelectionMode() ? 'payment.chooseMethodStartHint' : 'payment.chooseMethodHint';
        if (addressValue) addressValue.textContent = t(pendingAddressKey, t('payment.addressPending', 'Address will be generated after update'));
        if (copyAddressButton) copyAddressButton.disabled = true;
        renderQrCode('');
        if (paymentHintEl) paymentHintEl.textContent = t(pendingHintKey, t('payment.chooseMethodHint', 'Click update to refresh amount, address and QR code on this page.'));
    }

    function renderQrCode(address) {
        if (!qrcodeEl) return;

        qrcodeEl.innerHTML = '';
        if (!address) {
            const placeholder = document.createElement('div');
            placeholder.className = 'qr-placeholder';
            placeholder.textContent = t('payment.qrPending', 'QR code pending');
            qrcodeEl.appendChild(placeholder);
            return;
        }

        if (window.jQuery && window.jQuery.fn && window.jQuery.fn.qrcode) {
            const size = Math.max(160, Math.min(qrcodeEl.clientWidth || 210, 240));
            window.jQuery('#qrcode').qrcode({
                text: address,
                width: size,
                height: size,
                foreground: '#111827',
                background: '#ffffff',
                typeNumber: -1
            });
            return;
        }

        qrcodeEl.textContent = address;
    }

    function updatePayButton(hasSelected) {
        if (!payButton) return;

        const currentSelected = selectedMatchesCurrent();
        const locked = isSwitchLocked();
        payButton.disabled = !hasSelected || locked || currentSelected;
        payButton.classList.remove('is-loading');

        const label = payButton.querySelector('p');
        if (!label) return;

        if (locked && !currentSelected) {
            label.textContent = t('payment.switchLimitExhausted', 'Switch limit exhausted');
            return;
        }
        if (currentSelected) {
            label.textContent = t('payment.paymentReadyButton', 'Payment details ready');
            return;
        }
        label.textContent = isInitialSelectionMode()
            ? t('payment.startPayment', 'Start Payment')
            : t('payment.updatePaymentMethod', 'Update payment method');
    }

    function updateSelectorLockState() {
        const locked = !isWaitingOrder() || (!!currentPayment.address && remainingSwitchCount() <= 0);
        [currencySelectEl, networkSelectEl].forEach(function (el) {
            if (!el) return;
            if (locked) {
                el.classList.add('locked');
            } else {
                el.classList.remove('locked');
            }
        });
    }

    function createTransaction() {
        if (!selectedMethod || !payButton || payButton.disabled) return;

        payButton.disabled = true;
        payButton.classList.add('is-loading');

        fetch('/api/v1/pay/update-order', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                trade_id: tradeId,
                currency: selectedMethod.currency,
                network: selectedMethod.network
            })
        })
            .then(function (response) {
                return response.json();
            })
            .then(function (res) {
                if (res.status_code === 200 && res.data) {
                    const wasInitialSelection = isInitialSelectionMode();
                    applyPaymentPayload(res.data);
                    showMessage(wasInitialSelection ? t('payment.startSuccess', 'Payment started') : t('payment.updateSuccess', 'Payment method updated'));
                    return;
                }
                showMessage(res.message || t(isInitialSelectionMode() ? 'payment.createTransactionStartFailed' : 'payment.createTransactionFailed', 'Failed to update payment method'));
                updateUI();
            })
            .catch(function (err) {
                console.error(err);
                showMessage(t('payment.networkError', 'Network error'));
                updateUI();
            });
    }

    function applyPaymentPayload(payload) {
        currentPayment = payload || {};
        paymentConfig.current_payment = currentPayment;

        const expire = parseInt(currentPayment.expiration_time, 10);
        if (!isNaN(expire) && expire > 0) {
            totalSeconds = expire;
        }

        const currentMethod = findCurrentMethod();
        if (currentMethod) {
            selectedCurrency = currentMethod.currency;
            selectedNetwork = currentMethod.token_net_name;
            selectedMethod = currentMethod;
            selectOption(currencySelectEl, selectedCurrency);
            updateNetworkOptions(currentMethod);
        }

        updateContent();
        updateUI();
    }

    function showToast(msg) {
        let toast = document.querySelector('.cashier-toast');
        if (!toast) {
            toast = document.createElement('div');
            toast.className = 'cashier-toast';
            document.body.appendChild(toast);
        }

        toast.textContent = msg;
        toast.classList.add('show');
        if (toastTimer) clearTimeout(toastTimer);
        toastTimer = setTimeout(function () {
            toast.classList.remove('show');
        }, 2600);
    }

    function showMessage(msg) {
        if (window.layer && window.layer.msg) {
            window.layer.msg(msg);
            return;
        }
        showToast(msg);
    }

    function copyToClipboard(text, element, isButton) {
        if (!text) return;

        if (!navigator.clipboard || !navigator.clipboard.writeText) {
            fallbackCopy(text, element, isButton);
            return;
        }

        navigator.clipboard.writeText(text).then(function () {
            showCopySuccess(element, isButton);
        }).catch(function (err) {
            console.error('复制失败: ', err);
            fallbackCopy(text, element, isButton);
        });
    }

    function showCopySuccess(element, isButton) {
        const translate = window.i18next ? window.i18next.t.bind(window.i18next) : function (key) {
            const translations = {
                'payment.copied': '已复制!',
                'payment.copySuccess': '✓ 已复制!'
            };
            return translations[key] || key;
        };

        const originalText = element.textContent;
        if (isButton) {
            element.textContent = translate('payment.copied');
            element.classList.add('copied');
            setTimeout(function () {
                element.textContent = originalText;
                element.classList.remove('copied');
            }, 2000);
            return;
        }

        const originalColor = element.style.color;
        element.textContent = translate('payment.copySuccess');
        element.style.color = '#059669';
        setTimeout(function () {
            element.textContent = originalText;
            element.style.color = originalColor;
        }, 2000);
    }

    function fallbackCopy(text, element, isButton) {
        const textArea = document.createElement('textarea');
        textArea.value = text;
        textArea.style.position = 'fixed';
        textArea.style.opacity = '0';
        document.body.appendChild(textArea);
        textArea.select();
        try {
            document.execCommand('copy');
            showCopySuccess(element, isButton);
        } catch (err) {
            showMessage(t('payment.copyFailed', 'Copy failed, please copy manually'));
        }
        document.body.removeChild(textArea);
    }

    function copyAmount() {
        const active = selectedMatchesCurrent() ? currentPayment : selectedMethod;
        if (active && active.actual_amount) {
            copyToClipboard(active.actual_amount, amountValue, false);
        }
    }

    function copyAddress() {
        if (!selectedMatchesCurrent() || !currentPayment.address) return;
        copyToClipboard(currentPayment.address, copyAddressButton || addressValue, !!copyAddressButton);
    }

    function startCountdown() {
        if (countdownTimer) clearInterval(countdownTimer);

        let expireTime = parseInt(paymentConfig.expire, 10);
        if (isNaN(expireTime)) {
            expireTime = 0;
        }
        totalSeconds = expireTime;

        const minutesEl = document.getElementById('minutes');
        const secondsEl = document.getElementById('seconds');
        const countdownBanner = document.querySelector('.countdown-banner');

        function updateDisplay() {
            if (!minutesEl || !secondsEl) return;

            const minutes = Math.floor(Math.max(0, totalSeconds) / 60);
            const seconds = Math.max(0, totalSeconds) % 60;
            minutesEl.textContent = minutes.toString().padStart(2, '0');
            secondsEl.textContent = seconds.toString().padStart(2, '0');

            if (countdownBanner) {
                countdownBanner.classList.toggle('is-urgent', totalSeconds <= 300);
                countdownBanner.classList.toggle('is-critical', totalSeconds <= 60);
            }
        }

        updateDisplay();
        countdownTimer = setInterval(function () {
            totalSeconds--;
            updateDisplay();

            if (totalSeconds <= 0) {
                clearInterval(countdownTimer);
                if (statusCheckTimer) clearInterval(statusCheckTimer);
                showTimeoutMessage();
            }
        }, 1000);
    }

    function showTimeoutMessage() {
        if (document.querySelector('.status-overlay')) return;

        const overlay = document.createElement('div');
        overlay.className = 'status-overlay';

        const modal = document.createElement('div');
        modal.className = 'status-modal';
        modal.innerHTML = '<div class="status-icon">⏰</div>' +
            '<h3 class="status-title danger">' + t('payment.paymentTimeout', 'Payment Timeout') + '</h3>' +
            '<p class="status-text">' + t('payment.timeoutMessage', 'Sorry, the payment time has expired.<br>Please restart payment.') + '</p>' +
            '<button class="status-button" onclick="location.href=\'' + (paymentConfig.return_url || '/') + '\'">' +
            t('payment.returnToMerchant', 'Return to Merchant') + '</button>';

        overlay.appendChild(modal);
        document.body.appendChild(overlay);
    }

    function showWaitingConfirmation(data) {
        if (document.getElementById('waiting-overlay')) return;

        const overlay = document.createElement('div');
        overlay.id = 'waiting-overlay';
        overlay.className = 'status-overlay';

        const modal = document.createElement('div');
        modal.className = 'status-modal';
        modal.innerHTML = '<div class="status-spinner"></div>' +
            '<h3 class="status-title">' + t('payment.waitingConfirmation', 'Waiting for Confirmation') + '</h3>' +
            '<p class="status-text strong">' + t('payment.confirmationMessage', 'Payment detected, confirming...') + '</p>' +
            '<p class="status-note">' + t('payment.confirmationNote', 'Please wait patiently') + '</p>' +
            '<p class="status-note">' + t('payment.estimatedTime', 'Estimated time: 1-3 mins') + '</p>';

        overlay.appendChild(modal);
        document.body.appendChild(overlay);

        if (!statusCheckTimer) {
            statusCheckTimer = setInterval(checkPaymentStatus, 5000);
        }
    }

    function showSuccessMessage(data) {
        const waitingOverlay = document.getElementById('waiting-overlay');
        if (waitingOverlay) waitingOverlay.remove();
        if (document.querySelector('.success-overlay')) return;

        const overlay = document.createElement('div');
        overlay.className = 'status-overlay success-overlay';

        const modal = document.createElement('div');
        modal.className = 'status-modal';

        const icon = document.createElement('div');
        icon.className = 'status-icon success';
        icon.textContent = '✓';
        modal.appendChild(icon);

        const title = document.createElement('h3');
        title.className = 'status-title success';
        title.textContent = t('payment.paymentSuccess', 'Payment Success!');
        modal.appendChild(title);

        const text = document.createElement('p');
        text.className = 'status-text';
        text.textContent = t('payment.transactionHash', 'Transaction Hash:');
        modal.appendChild(text);

        const txHash = data.trade_hash || '';
        if (txHash) {
            const link = document.createElement('a');
            link.className = 'status-code status-code-link';
            link.href = 'https://www.oklink.com/search?q=' + encodeURIComponent(txHash);
            link.target = '_blank';
            link.rel = 'noopener noreferrer';
            link.title = t('payment.viewOnOklink', 'View on OKLink');
            link.textContent = txHash;
            modal.appendChild(link);
        } else {
            const code = document.createElement('code');
            code.className = 'status-code';
            code.textContent = t('payment.processing', 'Processing...');
            modal.appendChild(code);
        }

        const returnButton = document.createElement('button');
        returnButton.className = 'status-button success';
        returnButton.type = 'button';
        returnButton.textContent = t('payment.returnToMerchant', 'Return to Merchant');
        returnButton.addEventListener('click', function () {
            location.href = data.return_url || '/';
        });
        modal.appendChild(returnButton);

        overlay.appendChild(modal);
        document.body.appendChild(overlay);

        if (countdownTimer) clearInterval(countdownTimer);
        if (statusCheckTimer) clearInterval(statusCheckTimer);
    }

    function checkPaymentStatus() {
        if (!tradeId) return;

        fetch('/pay/check-status/' + tradeId)
            .then(function (r) {
                return r.json();
            })
            .then(function (data) {
                if (data.status === 2) {
                    showSuccessMessage(data);
                } else if (data.status === 3) {
                    showTimeoutMessage();
                } else if (data.status === 5) {
                    showWaitingConfirmation(data);
                }
            })
            .catch(function (err) {
                console.error('Check status error:', err);
            });
    }

    function init(config) {
        paymentConfig = config || {};
        currentPayment = paymentConfig.current_payment || {};
        tradeId = paymentConfig.trade_id;

        currencySelectEl = document.getElementById('currency-select');
        networkSelectEl = document.getElementById('network-select');
        payButton = document.querySelector('.pay-button');
        amountValue = document.getElementById('payAmount');
        addressValue = document.getElementById('walletAddress');
        copyAddressButton = document.getElementById('copyAddressButton');
        qrcodeEl = document.getElementById('qrcode');
        paymentHintEl = document.getElementById('paymentHint');
        exchangeRateEl = document.getElementById('exchangeRate');
        selectedCurrencyEl = document.getElementById('selectedCurrency');
        selectedNetworkEl = document.getElementById('selectedNetwork');
        switchCounterEl = document.getElementById('switchCounter');

        document.addEventListener('click', function () {
            document.querySelectorAll('.coin-select').forEach(function (el) {
                el.classList.remove('active');
            });
        });

        setupDropdown(currencySelectEl);
        setupDropdown(networkSelectEl);

        if (payButton) {
            payButton.addEventListener('click', createTransaction);
        }
        if (copyAddressButton) {
            copyAddressButton.addEventListener('click', copyAddress);
        }

        initI18n().then(function () {
            if (paymentConfig.network && Array.isArray(paymentConfig.network) && paymentConfig.network.length > 0) {
                paymentMethods = paymentConfig.network;
                initSelection();
            } else {
                showMessage(t('payment.loadPaymentNetworkFailed', 'Failed to load payment network, please check if the wallet has been added'));
                updateUI();
            }

            startCountdown();
            statusCheckTimer = setInterval(checkPaymentStatus, 5000);
            checkPaymentStatus();
        });
    }

    window.copyAmount = copyAmount;
    window.copyAddress = copyAddress;
    window.changeLanguage = changeLanguage;
    window.Payment = {
        init: init
    };
})();
