a requisição envia os dados da simulação pro servidor => servidor verifica o número de parcelas. se for maior que 06, calcula os juros. se não for, só devolve o valor total a ser pago com o número de parcelas. essa resposta é armazenada no redis pro tempo limitado. dentro desse período, o cliente pode confirmar a compra. se for confirmada, os dados da compra são enviados ao postgresql.