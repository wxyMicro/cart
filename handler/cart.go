package handler

import (
	"context"

	"github.com/wxyMicro/cart/domain/model"
	"github.com/wxyMicro/cart/domain/service"
	cart "github.com/wxyMicro/cart/proto/cart"

	common "github.com/wxyMicro/common"
)

type Cart struct {
	CartDataService service.ICartDataService
}

// AddCart 添加购物车
func (c *Cart) AddCart(ctx context.Context, request *cart.CartInfo, response *cart.ResponseAdd) (err error) {
	cartData := &model.Cart{}
	_ = common.SwapTo(request, cartData)
	response.CartId, err = c.CartDataService.AddCart(cartData)
	return err
}

// CleanCart 根据用户id清空用户购物车
func (c *Cart) CleanCart(ctx context.Context, request *cart.Clean, response *cart.Response) error {
	err := c.CartDataService.CleanCart(request.UserId)
	if err != nil {
		return err
	}
	response.Msg = "购物车清空成功"
	return nil
}

// Incr 添加购物车商品数量
func (c *Cart) Incr(ctx context.Context, request *cart.Item, response *cart.Response) error {
	err := c.CartDataService.IncrNum(request.Id, request.ChangeNum)
	if err != nil {
		return err
	}
	response.Msg = "购物车添加成功"
	return nil
}

// Decr 减少购物车商品数量
func (c *Cart) Decr(ctx context.Context, request *cart.Item, response *cart.Response) error {
	if err := c.CartDataService.DecrNum(request.Id, request.ChangeNum); err != nil {
		return err
	}
	response.Msg = "购物车减少成功"
	return nil
}

// DeleteItemByID 根据购物车id删除购物车
func (c *Cart) DeleteItemByID(ctx context.Context, request *cart.CartID, response *cart.Response) error {
	err := c.CartDataService.DeleteCart(request.Id)
	if err != nil {
		return err
	}
	response.Msg = "购物车删除成功"
	return nil
}

//  查寻用户所有的购物车信息
func (c *Cart) GetAll(ctx context.Context, request *cart.CartFindAll, response *cart.CartAll) error {
	allCart, err := c.CartDataService.FindAllCart(request.UserId)
	if err != nil {
		return err
	}
	for _, v := range allCart {
		cartInfo := &cart.CartInfo{}
		err := common.SwapTo(v, cartInfo)
		if err != nil {
			return err
		}
		response.CartInfo = append(response.CartInfo, cartInfo)
	}
	return nil
}
