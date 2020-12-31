<template>
  <Layout :isLoading="isLoading" backgroundColor="#dee2e6">
    <template #default>
      <div class="container text-left">
        <div class="row">
          <div class="col">
            <div class="row">
              <div class="col mb-3 mb-sm-5">
                <div class="quantity-readout p-3 font-weight-bold">{{ cartQuantityReadout }}</div>
              </div>
            </div>

            <ul class="cart-items">
              <CartItem
                v-for="item in cart.items"
                :key="item.product_id"
                :product_id="item.product_id"
                :quantity="item.quantity"
                class="mb-4"
              ></CartItem>
            </ul>
          </div>


          <div v-if="cart.items.length > 0" class="summary-container col-lg-auto">
            <div class="summary p-4">
              <div class="summary-quantity">{{ summaryQuantityReadout }}</div>
              <div class="summary-total mb-2 font-weight-bold">Your Total: {{ formattedCartTotal }}</div>
              <router-link to="/checkout" class="checkout-btn mb-3 btn btn-outline-dark btn-block btn-lg"
                >Checkout
              </router-link>
              <button v-if="pinpointEnabled && user" v-on:click="triggerAbandonedCartEmail" class="abandoned-cart-btn btn btn-primary btn-block btn-lg">
                Trigger Abandoned Cart email
              </button>
            </div>
          </div>
        </div>

        <!-- <RecommendedProductsSection :feature="feature" :recommendedProducts="featuredProducts">
        <template #heading>
          Featured products
        </template>
        </RecommendedProductsSection> -->

        <RecommendedProductsSection
          :explainRecommended="explainRecommended"
          :recommendedProducts="relatedProducts"
          :feature="feature"
        >
          <template #heading>Compare similar items</template>
        </RecommendedProductsSection>
      </div>
    </template>
  </Layout>
</template>

<script>
import swal from 'sweetalert';

import { mapState, mapActions, mapGetters } from 'vuex';
import { RepositoryFactory } from '@/repositories/RepositoryFactory';


import { AnalyticsHandler } from '@/analytics/AnalyticsHandler'

import CartItem from './components/CartItem.vue';
import Layout from '@/components/Layout/Layout';
import RecommendedProductsSection from '@/components//RecommendedProductsSection/RecommendedProductsSection';
const RecommendationsRepository = RepositoryFactory.get('recommendations');
const MAX_RECOMMENDATIONS = 9;
const EXPERIMENT_FEATURE = 'home_product_recs';


export default {
  name: 'Cart',
  components: {
    Layout,
    CartItem,
    RecommendedProductsSection
  },
  props: {
  },
  data () {
    return {  
      pinpointEnabled : process.env.VUE_APP_PINPOINT_APP_ID,
      feature: EXPERIMENT_FEATURE,
      featuredProducts: null,
      quantity: 1,
      relatedProducts: null,
      explainRecommended: null,
    }
  },
  created() {
        AnalyticsHandler.cartViewed(
          this.user,
          this.cart,
          this.cartQuantity,
          this.cartTotal,
        );
  },
  computed: {
    ...mapState({ cart: (state) => state.cart.cart, user: state => state.user }),
    ...mapGetters(['cartQuantity', 'cartTotal', 'formattedCartTotal']),
    isLoading() {
      return !this.cart;
    },
   
        

    cartQuantityReadout() {
      if (this.cartQuantity === null) return null;

      return `(${this.cartQuantity}) ${this.cartQuantity === 1 ? 'item' : 'items'} in your cart shopping cart`;

    },
    summaryQuantityReadout() {
      if (this.cartQuantity === null) return null;

      return `Summary (${this.cartQuantity}) ${this.cartQuantity === 1 ? 'item' : 'items'}`;
    },
  },
  
 watch: {
    $route: {
      immediate: true,
      handler() {
        this.fetchData();
      },
    },
  },
  methods: {
     ...mapActions(['addToCart']),
    resetQuantity() {
      this.quantity = 1;
    },
    async addProductToCart() {
      await this.addToCart({
        product: this.product,
        quantity: this.quantity,
        feature: this.$route.query.feature,
        exp: this.$route.query.exp,
      });

      this.renderAddedToCartConfirmation();

      this.resetQuantity();
    },
    async triggerAbandonedCartEmail () {
      if (this.cart && this.cart.items.length > 0 ){
      const cartItem = await this.cart.items[0].product_id

      AnalyticsHandler.recordAbanonedCartEvent(this.user,this.cart,cartItem)
      }
      else{
        console.error("No items to export")
      }
    },
    async fetchData() {

      // reset in order to trigger recalculation in carousel - carousel UI breaks without this
      this.relatedProducts = null;
      this.getRelatedProducts();

    },
    async getRelatedProducts() {
      const response = await RecommendationsRepository.getRelatedProducts(
        this.personalizeUserID ?? '',
        this.cart.items[0].product_id,
        MAX_RECOMMENDATIONS,
        EXPERIMENT_FEATURE,
      );

      if (response.headers) {
        const experimentName = response.headers['x-experiment-name'];
        const personalizeRecipe = response.headers['x-personalize-recipe'];

        if (experimentName || personalizeRecipe) {
          const explanation = experimentName
            ? `Active experiment: ${experimentName}`
            : `Personalize recipe: ${personalizeRecipe}`;

          this.explainRecommended = {
            activeExperiment: !!experimentName,
            personalized: !!personalizeRecipe,
            explanation,
          };
        }
      }

      this.relatedProducts = response.data;

      if (this.relatedProducts.length > 0 && 'experiment' in this.relatedProducts[0]) {
        AnalyticsHandler.identifyExperiment(this.user, this.relatedProducts[0].experiment);
      }
    },
    
    renderAddedToCartConfirmation() {
      swal({
        title: 'Added to Cart',
        icon: 'success',
        buttons: {
          cancel: 'Continue Shopping',
          cart: 'View Cart',
        },
      }).then((value) => {
        switch (value) {
          case 'cancel':
            break;
          case 'cart':
            this.$router.push('/cart');
        }
      });
    },
    
  },
};
</script>

<style scoped>
.quantity-readout {
  border-radius: 4px;
  background: var(--grey-200);
  font-size: 1.15rem;
}

.cart-items {
  list-style-type: none;
  padding: 0;
}

.summary {
  border: 1px solid var(--grey-600);
  border-radius: 2px;
}

.summary-total {
  font-size: 1.15rem;
}

.summary-quantity {
  font-size: 1.15rem;
}

.abandoned-cart-btn {
  background: var(--blue-500);
  border-color: var(--blue-500);
  font-size: 1rem;
}

.abandoned-cart-btn:hover,
.abandoned-cart-btn:focus {
  background: var(--blue-600);
  border-color: var(--blue-600);
}

.checkout-btn {
  border-color: var(--grey-900);
  border-width: 2px;
  font-size: 1rem;
}

.checkout-btn:hover,
.checkout-btn:focus {
  background: var(--grey-900);
}

@media (min-width: 768px) {
  .quantity-readout {
    font-size: 1.75rem;
  }

  .summary-total {
    font-size: 1.5rem;
  }

  .summary-quantity {
    font-size: 1.5rem;
  }

  .abandoned-cart-btn,
  .checkout-btn {
    font-size: 1.25rem;
  }
}

@media (min-width: 992px) {
  .summary-container {
    min-width: 350px;
  }

  .summary {
    position: sticky;
    top: 120px;
  }
}
</style>
