# What is CQRS

In his book "[Object Oriented Software Construction][meyerbook]," 
Betrand Meyer introduced the term "[Command Query Separation][cqsterm]" 
to describe the principle that an object's methods should be either 
commands, or queries. A query returns data and does not alter the state 
of the object; a command changes the state of an object but does not 
return any data. The benefit is that you have a better understanding 
what does, and what does not, change the state in your system. 

# Read and write sides
When you apply the CQRS pattern, you can separate the read and write sides in this portion of the system.

![Figure 1][fig1]
**A possible architectural implementation of the CQRS pattern**

# Why should I use CQRS?

The most common business benefits that you might gain from applying the CQRS pattern are enhanced scalability, the simplification of a complex aspect of your domain, increased flexibility of your solution, and greater adaptability to changing business requirements.

## Scalability 

In distributed system, the number of reads vastly exceeds the number of writes, so your scalability requirements will be different for each side. By separating the read-side and the write-side into separate models, you now have the ability to scale each one of them independently.

## Reduced complexity
=))
https://vladikk.com/2017/03/20/tackling-complexity-in-cqrs/ 

## Flexibility

It becomes much easier to make changes on the read-side, such as adding a new query to support a new report screen in the UI, when you can be confident that you won't have any impact on the behavior of the business logic. On the write-side, having a model that concerns itself solely with the core business logic in the domain means that you have a simpler model to deal with than a model that includes read logic as well.

[meyerbook]:      http://www.amazon.com/gp/product/0136291554
[cqsterm]:        http://martinfowler.com/bliki/CommandQuerySeparation.html
[fig1]:           images/Reference_02_Arch_01.png?raw=true